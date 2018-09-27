package config

import (
	"bytes"
	"os"
	"testing"

	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
)

const (
	testENV = "DEBUG"
)

func unsetEnv() {
	os.Unsetenv(serviceUUID)
	os.Unsetenv(serviceHost)
	os.Unsetenv(ns)
	os.Unsetenv(nsListener)
	os.Unsetenv(nsPublisher)
	os.Unsetenv(serviceName)
	os.Unsetenv(serviceProcess)
	os.Unsetenv("MQTT_LISTENER_URL")
	os.Unsetenv("MQTT_PUBLISHER_URL")
}

func TestLoad(t *testing.T) {
	defer unsetEnv()
	logger.Log = &logrus.Logger{}
	testCases := []struct {
		name       string
		needErr    bool
		configPath string
		setEnv     bool
	}{
		{"Test Load with bad ConfigPath", true, "_test.json", false},
		{"Test Load with empty ConfigName", true, "", false},
		{"Test Load with empty ConfigName", false, "", true},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unsetEnv()
			ConfigPath = tc.configPath
			if tc.setEnv {
				os.Setenv(serviceName, "test")
				os.Setenv("MQTT_LISTENER_URL", "tcp://golang.org:443")
				os.Setenv("MQTT_PUBLISHER_URL", "tcp://golang.org:443")
			}
			err = Load()
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func TestInitJSON(t *testing.T) {
	logger.Log = &logrus.Logger{}
	config := new(Configuration)
	testCases := []struct {
		name       string
		needErr    bool
		configPath string
	}{
		{"Test initJSON with empty ConfigPath", false, ""},
		{"Test initJSON with empty ConfigPath", false, "qwertyuiopasdfghjkl"},
		{"Test initJSON with empty ConfigPath", false, "qwertyuiopasdfghjkl"},
		{"Test initJSON with empty ConfigPath", true, "_test.json"},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = initJSON(config, tc.configPath)
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func TestInitEnv(t *testing.T) {
	defer os.Unsetenv(testENV)
	logger.Log = &logrus.Logger{}
	config := new(Configuration)
	testCases := []struct {
		name    string
		needErr bool
		env     map[string]string
	}{
		{"Test initEnv with empty environment", false, nil},
		{"Test initEnv with correct environment", false, map[string]string{testENV: "true"}},
		{"Test initEnv with not correct environment", true, map[string]string{testENV: "111"}},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.env {
				os.Setenv(key, val)
			}
			err = initEnv(config)
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

//
func TestProcessConfig(t *testing.T) {
	defer unsetEnv()
	logger.Log = &logrus.Logger{}
	config := new(Configuration)
	testCases := []struct {
		name       string
		needErr    bool
		configPath string
		env        map[string]string
	}{
		{"Test processConfig with empty configPath", false, "", nil},
		{"Test processConfig with not correct package.json file", true, "_test.json", nil},
		{"Test processConfig with not correct environment", true, "", map[string]string{testENV: "111"}},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ConfigPath = tc.configPath
			for key, val := range tc.env {
				os.Setenv(key, val)
			}
			err = processConfig(config)
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestConfig_setServiceProcessor(t *testing.T) {
	defer unsetEnv()
	config := new(Configuration)
	testCases := []struct {
		name   string
		setEnv bool
	}{
		{"Test setServiceProcessor with unset  SERVICE_PROCESSOR", false},
		{"Test setServiceProcessor with unset  SERVICE_PROCESSOR", true},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				os.Setenv(serviceProcess, "")
			}
			err = config.setServiceProcessor()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestConfig_setNamespace(t *testing.T) {
	defer unsetEnv()
	config := new(Configuration)
	testCases := []struct {
		name    string
		needErr bool
		setEnv  bool
		ns      string
		nsList  string
		nsPub   string
	}{
		{"Test setNamespace with unset  NAMESPACE", false, false, "", "", ""},
		{"Test setNamespace with set not correct  NAMESPACE", true, false, string([]byte{0}), "", ""},
		{"Test setNamespace with set not correct  NAMESPACE_LISTENER", true, true, string([]byte{0}), "", ""},
		{"Test setNamespace with set not correct  NAMESPACE_PUBLISHER", true, true, string([]byte{0}), "test", ""},
		{"Test setNamespace with correct  environment", false, true, string([]byte{0}), "test", "test"},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Unsetenv(ns)
			if tc.setEnv {
				os.Setenv(ns, "")
			}
			config.Namespace = tc.ns
			config.NamespaceListener = tc.nsList
			err = config.setNamespace()
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}

		})
	}
}

//
func TestConfig_setName(t *testing.T) {
	defer unsetEnv()
	config := new(Configuration)
	testCases := []struct {
		name     string
		needErr  bool
		setEnv   bool
		confName string
	}{
		{"Test setName with unset Configuration.NAME", true, false, ""},
		{"Test setName with set SERVICE_NAME", false, true, "test"},
		{"Test setName with unset good SERVICE_NAME", false, false, "test"},
		{"Test setName with unset bad SERVICE_NAME", true, false, string([]byte{0})},
	}
	for _, tc := range testCases {
		var err error
		t.Run(tc.name, func(t *testing.T) {
			os.Unsetenv(serviceName)
			if tc.setEnv {
				os.Setenv(serviceName, tc.confName)
			}
			config.Name = tc.confName
			err = config.setName()
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestConfig_setURL(t *testing.T) {
	config := new(Configuration)
	testCases := []struct {
		name    string
		needErr bool
		pubURL  string
		listURL string
	}{
		{"Test setURL with notExisted MQTT_PUBLISHER_URL", true, "tcp://:", ""},
		{"Test setURL with bad MQTT_PUBLISHER_URL", true, ":", ""},
		{"Test setURL with bad MQTT_LISTENER_URL", true, "tcp://golang.org:443", ":"},
		{"Test setURL with good MQTT_LISTENER_URL", false, "tcp://golang.org:443", "tcp://golang.org:443"},
	}
	for _, tc := range testCases {
		var err error
		t.Run(tc.name, func(t *testing.T) {
			config.MQTTPublisherURL = tc.pubURL
			config.MQTTListenerURL = tc.listURL
			err = config.setURL()
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestConfig_checkConfig(t *testing.T) {
	defer unsetEnv()
	var buf bytes.Buffer
	logger.Log = logrus.New()
	logger.Log.Out = &buf
	config := &Configuration{
		MQTTListenerURL:    "tcp://golang.org:443",
		MQTTPublisherURL:   "tcp://golang.org:443",
		Namespace:          "test",
		NamespaceListener:  "test",
		NamespacePublisher: "test",
		ServiceProcessor:   "test",
		Topic:              "test",
	}
	testCases := []struct {
		name       string
		needErr    bool
		configName string
		debugLevel bool
	}{
		{"Test checkConfig with empty SERVICE_NAME", true, "", false},
		{"Test checkConfig with set SERVICE_NAME", false, "test", true},
	}
	for _, tc := range testCases {
		var err error
		t.Run(tc.name, func(t *testing.T) {
			config.Name = tc.configName
			config.Debug = tc.debugLevel
			err = config.checkConfig()
			if tc.needErr {
				if err == nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestConfig_CheckMQTT(t *testing.T) {
	logger.Log = &logrus.Logger{}
	credoP := Credentials{
		Password: "test2",
		UserName: "test2",
	}
	config := new(Configuration)
	if !config.checkMQTT() {
		t.Error("unexpected result, expected true")
	}
	config.PubCredo = credoP
	if config.checkMQTT() {
		t.Error("unexpected result, expected false")
	}
}

func TestConfig_setTopic(t *testing.T) {
	logger.Log = &logrus.Logger{}
	config := new(Configuration)
	config.setTopic()
	if config.Topic != "" {
		t.Errorf("expected empty topic: %s", config.Topic)
	}
	SubscriptionsPath = "_test.txt"
	config.setTopic()
	if config.Topic == "" {
		t.Error("expected not empty topic")
	}
}

func TestGetCredo(t *testing.T) {
	logger.Log = &logrus.Logger{}
	credo := getCredo("_test.txt")
	if credo.UserName != "" || credo.Password != "" {
		t.Errorf("unexpected result: %v", credo)
	}
	credo = getCredo("_test.json")
	if credo.UserName == "" || credo.Password == "" {
		t.Errorf("unexpected result: %v", credo)
	}
}
