package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"mqtt-adapter/src/logger"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	ns             = "NAMESPACE"
	nsListener     = "NAMESPACE_LISTENER"
	nsPublisher    = "NAMESPACE_PUBLISHER"
	serviceHost    = "SERVICE_HOST"
	serviceUUID    = "SERVICE_UUID"
	serviceName    = "SERVICE_NAME"
	serviceProcess = "SERVICE_PROCESSOR"

	timeOut = time.Second * 10
	tcp     = "tcp"
)

var (
	// ConfigPath represents default path to config file
	ConfigPath = "path/to/package.json"
	// SubscriptionsPath represents default path to subscriptions.txt file
	SubscriptionsPath = "path/to/subscriptions.txt"
	// Config is a container for Configuration information
	Config *Configuration
	// PubCredoPath is a path to Publisher secret
	PubCredoPath = "path/to/secrets/mqtt_publisher.json"
	// ListCredoPath is a path to Listener secret
	ListCredoPath = "path/to/secrets/mqtt_listener.json"
)

// Credentials is a container for MQTT credentials
type Credentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Configuration represents Configuration options
type Configuration struct {
	Name               string `envconfig:"SERVICE_NAME"          json:"name"`
	UUID               string
	Host               string
	MQTTListenerURL    string `envconfig:"MQTT_LISTENER_URL"     default:"tcp://mqtt:1883"`
	MQTTPublisherURL   string `envconfig:"MQTT_PUBLISHER_URL"    default:"tcp://mqtt:1883"`
	Namespace          string `envconfig:"NAMESPACE"             default:"default"`
	NamespaceListener  string `envconfig:"NAMESPACE_LISTENER"`
	NamespacePublisher string `envconfig:"NAMESPACE_PUBLISHER"`
	ServiceProcessor   string `envconfig:"SERVICE_PROCESSOR"     default:"./service-processor/processor"`
	Topic              string
	ListCredo          Credentials
	PubCredo           Credentials
	Debug              bool   `envconfig:"DEBUG"`
	Bridge             bool   `envconfig:"BRIDGE"`
	Same               bool
}

func (c *Configuration) setTopic() error {
	logger.Log.Infof("Trying to read file %q ... ", SubscriptionsPath)
	subscription, err := ioutil.ReadFile(SubscriptionsPath)
	if err != nil {
		msg := fmt.Sprintf("Reading Subscriptions.txt (%s) failed (err: %v).\n", SubscriptionsPath, err)
		logger.Log.Warn(msg)
		return nil
	}
	c.Topic = strings.TrimSpace(string(subscription))
	logger.Log.Infof("Topic subscribed %s", c.Topic)
	return nil
}

// setSecrets reads secrets from json files
func (c *Configuration) setSecrets() error {
	c.PubCredo = getCredo(PubCredoPath)
	c.ListCredo = getCredo(ListCredoPath)
	return nil
}

func (c *Configuration) checkOnSame() error {
	c.Same = c.checkMQTT()
	return nil
}

// checkMQTT log message about MQTT publisher and listener servers
func (c *Configuration) checkMQTT() bool {
	if c.MQTTListenerURL != c.MQTTPublisherURL ||
		c.ListCredo.UserName != c.PubCredo.UserName ||
		c.ListCredo.Password != c.PubCredo.Password {
		return false
	}
	logger.Log.Debugln("MQTT connection: listener and publisher are equal")
	return true
}

// getCredo reads credentials from MQTT secret files
func getCredo(filePath string) Credentials {
	credo := new(Credentials)
	secretFile, err := ioutil.ReadFile(filePath)
	logger.Log.Debugf("Trying to read file %q ... ", filePath)
	if err != nil {
		logger.Log.Warnf("Reading MQTT Secret configuration from JSON (%s) failed (err: %v). SKIPPED.\n", filePath, err)
		return *credo
	}
	logger.Log.Debugln("OK")
	reader := bytes.NewBuffer(secretFile)
	json.NewDecoder(reader).Decode(credo)
	return *credo
}

// setServiceProcessor checks if specified process exists
func (c *Configuration) setServiceProcessor() (err error) {
	_, found := syscall.Getenv(serviceProcess)
	if found {
		return nil
	}
	return syscall.Setenv(serviceProcess, c.ServiceProcessor)
}

// setUUID generates random UUID and sets it to system environment SERVICE_UUID
func (c *Configuration) setUUID() (err error) {
	var u uuid.UUID
	u, _ = uuid.NewV4()
	c.UUID = u.String()
	return syscall.Setenv(serviceUUID, u.String())
}

// setHostName retrieves host name of machine and sets it to system environment SERVICE_HOST
func (c *Configuration) setHostName() (err error) {
	var hostName string
	hostName, _ = os.Hostname()
	c.Host = hostName
	return syscall.Setenv(serviceHost, hostName)
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

// setURL checks if MQTT_LISTENER_URL and MQTT_PUBLISHER_URL are valid URL
func (c *Configuration) setURL() (err error) {
	if err = checkTCPConnection(c.MQTTPublisherURL); err != nil {
		return err
	}
	return checkTCPConnection(c.MQTTListenerURL)
}

// checkTCPConnection tries to connect to the address on the TCP network.
func checkTCPConnection(path string) (err error) {
	if _, err = url.Parse(path); err != nil {
		return err
	}
	dialer := net.Dialer{Timeout: timeOut}
	con, err := dialer.Dial(tcp, strings.TrimPrefix(path, "tcp://"))
	if err != nil {
		return err
	}
	con.Close()
	return nil
}

// Load initializes Config
func Load() (err error) {
	config := new(Configuration)
	if err = processConfig(config); err != nil {
		return err
	}
	err = config.checkConfig()
	if err != nil {
		return err
	}
	Config = config
	return nil
}

// checkConfig checks if config information is valid
func (c *Configuration) checkConfig() (err error) {
	configSetters := setList{
		{"setUUID", c.setUUID},
		{"setHostName", c.setHostName},
		{"setNamespace", c.setNamespace},
		{"setName", c.setName},
		{"setURL", c.setURL},
		{"setServiceProcess", c.setServiceProcessor},
		{"setTopic", c.setTopic},
		{"setSecrets", c.setSecrets},
		{"check on Same", c.checkOnSame},
	}
	if c.Debug {
		logger.Log.SetLevel(logrus.DebugLevel)
	}
	for _, setter := range configSetters {
		err := setter.set()
		if err != nil {
			logger.Log.Errorln(err)
			return fmt.Errorf("failed to execute Config.%s: %v", setter.name, err)
		}
	}
	return nil
}

// setList is a collection of Set() methods
type setList []struct {
	name string
	set  func() error
}

// processConfig try to load Configuration from Environment, File or by Default
func processConfig(config *Configuration) (err error) {
	initDefault(config)

	err = initJSON(config, ConfigPath)
	if err != nil {
		return err
	}
	return initEnv(config)
}

// initDefault initializes Configuration by default value
func initDefault(config *Configuration) (err error) {
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
		}
	}
	return
}

// initJSON initializes Configuration from specified JSON file
func initJSON(config *Configuration, filePath string) (err error) {
	if len(filePath) == 0 {
		return
	}
	logger.Log.Infof("Trying to read file %q ... ", filePath)
	// log.Printf("Trying to read file %q ... ", filePath)
	configFileContents, err := ioutil.ReadFile(filePath)

	if err != nil {
		logger.Log.Warnf("Reading Configuration from JSON (%s) failed (err: %v). SKIPPED.\n", filePath, err)
		return nil
	}
	logger.Log.Infoln("OK")
	reader := bytes.NewBuffer(configFileContents)
	return json.NewDecoder(reader).Decode(config)
}

// initEnv initializes Configuration from system environment
func initEnv(config *Configuration) (err error) {
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
