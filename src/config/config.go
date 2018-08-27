package config

import (
	"reflect"
	"fmt"
	"strconv"
	"io/ioutil"
	"log"
	"bytes"
	"encoding/json"
	"syscall"
	"os"
	"net/url"

	"github.com/satori/go.uuid"
)

const (
	ns          = "NAMESPACE"
	nsListener  = "NAMESPACE_LISTENER"
	nsPublisher = "NAMESPACE_PUBLISHER"
	serviceHost = "SERVICE_HOST"
	serviceUUID = "SERVICE_UUID"
	serviceName = "SERVICE_NAME"
)

var (
	Config Configuration
	// configPath represents default path to config file
	ConfigPath        = "package.json"
	SubscriptionsPath = "./service-processor/subscriptions.txt"
)

// Configuration represents configuration options
type Configuration struct {
	Name               string `envconfig:"SERVICE_NAME"          json:"name"`
	UUID               string
	Host               string
	Debug              bool   `envconfig:"DEBUG"`
	Bridge             bool   `envconfig:"BRIDGE"`
	MQTTListenerURL    string `envconfig:"MQTT_LISTENER_URL"     default:"tcp://mqtt:1883"`
	MQTTPublisherURL   string `envconfig:"MQTT_PUBLISHER_URL"    default:"tcp://mqtt:1883"`
	Namespace          string `envconfig:"NAMESPACE"             default:"default"`
	NamespaceListener  string `envconfig:"NAMESPACE_LISTENER"`
	NamespacePublisher string `envconfig:"NAMESPACE_PUBLISHER"`
	ServiceProcessor   string `envconfig:"SERVICE_PROCESSOR"     default:"./service-processor/processor"`
	LogLevel           string `envconfig:"LOG_LEVEL"             default:"error"`
}

// setUUID generates random UUID and sets it to system environment SERVICE_UUID
func (c *Configuration) setUUID() (err error) {
	var u uuid.UUID
	u, err = uuid.NewV4()
	if err != nil {
		return err
	}
	c.UUID = u.String()
	if err = syscall.Setenv(serviceUUID, u.String()); err != nil {
		return err
	}
	return nil
}

// setHostName retrieves host name of machine and sets it to system environment SERVICE_HOST
func (c *Configuration) setHostName() (err error) {
	var hostName string
	hostName, err = os.Hostname()
	if err != nil {
		return err
	}
	c.Host = hostName
	if err = syscall.Setenv(serviceHost, hostName); err != nil {
		return err
	}
	return nil
}

// setNamespace checks if system environments NAMESPACE, NAMESPACE_LISTENER and
// NAMESPACE_PUBLISHER were set. If no, method sets them to $NAMESPACE value
func (c *Configuration) setNamespace() (err error) {
	_, found := syscall.Getenv(ns)
	if !found {
		if err = syscall.Setenv(ns, c.Namespace); err != nil {
			return err
		}
	}
	if c.NamespaceListener == "" {
		c.NamespaceListener = c.Namespace
		if err = syscall.Setenv(nsListener, c.NamespaceListener); err != nil {
			return err
		}
	}
	if c.NamespacePublisher == "" {
		c.NamespacePublisher = c.Namespace
		if err = syscall.Setenv(nsPublisher, c.NamespacePublisher); err != nil {
			return err
		}
	}
	return nil
}

// setName checks if an environment SERVICE_NAME was set,
// if no, it would set to value from package.json file
func (c *Configuration) setName() (err error) {
	if c.Name == "" {
		return fmt.Errorf("SERVICE_NAME wasn't set, not from file, not from environment")
	}
	_, found := syscall.Getenv(serviceName)
	if found {
		return nil
	}
	return syscall.Setenv(serviceName, c.Name)
}

// checkURL checks if MQTT_LISTENER_URL and MQTT_PUBLISHER_URL are valid URL
func (c *Configuration) checkURL() (err error) {
	if _, err = url.Parse(c.MQTTPublisherURL); err != nil {
		return err
	}
	if _, err = url.Parse(c.MQTTListenerURL); err != nil {
		return err
	}
	return nil
}

// LoadConfig initializes Config
func LoadConfig() error {
	if err := processConfig(); err != nil {
		return err
	}

	for _, setter := range configSetters {
		err := setter.set()
		if err != nil {
			return fmt.Errorf("failed to execute Config.%s: %v", setter.name, err)
		}
	}
	return nil
}

// setList is a collection of Set() functions
type setList []struct {
	name string
	set  func() error
}

var configSetters = setList{
	{"setUUID", Config.setUUID},
	{"setHostName", Config.setHostName},
	{"setNamespace", Config.setNamespace},
	{"setName", Config.setName},
	{"checkURL", Config.checkURL},
}

func processConfig() (err error) {
	err = initConfigWithDefaultValues(&Config)
	if err != nil {
		return err
	}
	err = mergeJSONConfig(&Config, ConfigPath)
	if err != nil {
		return err
	}
	err = mergeEnvConfig(&Config)
	if err != nil {
		return err
	}
	return nil
}

func initConfigWithDefaultValues(config *Configuration) (err error) {
	configElements := reflect.TypeOf(config).Elem()
	for i := 0; i < configElements.NumField(); i++ {
		defaultValue, hasDefaultTag := configElements.Field(i).Tag.Lookup("default")
		if !hasDefaultTag {
			continue
		}
		field := reflect.ValueOf(config).Elem().Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(defaultValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(defaultValue)
			if err != nil {
				log.Printf("Cannot parse field %s value %v as bool type", field.Type().Name(), defaultValue)
			}
			field.SetBool(boolValue)
		}
	}
	return
}
func mergeJSONConfig(config *Configuration, filePath string) (err error) {
	if len(filePath) == 0 {
		return
	}
	log.Printf("Trying to read file %q ... ", filePath)
	configFileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed (err: %v). SKIPPED.\n", filePath, err)
		return nil
	}
	log.Println("OK")
	reader := bytes.NewBuffer(configFileContents)
	return json.NewDecoder(reader).Decode(config)
}

func mergeEnvConfig(config *Configuration) (err error) {
	configElements := reflect.TypeOf(config).Elem()
	for i := 0; i < configElements.NumField(); i++ {
		envKey, hasEnvconfigTag := configElements.Field(i).Tag.Lookup("envconfig")
		if !hasEnvconfigTag {
			continue
		}
		envValue, found := syscall.Getenv(envKey)
		if !found {
			continue
		}
		structFieldName := configElements.Field(i).Name
		envField := reflect.ValueOf(config).Elem().FieldByName(structFieldName)
		switch envField.Kind() {
		case reflect.String:
			envField.SetString(envValue)
		case reflect.Bool:
			boolEnvValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("cannot parse environment %s=%v as bool type", envKey, envValue)
			}
			envField.SetBool(boolEnvValue)
		}
	}
	return
}
