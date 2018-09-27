package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"mqtt-adapter/src/config"
	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
	"github.com/surgemq/surgemq/service"
)

func TestSetConfig(t *testing.T) {
	configTest := config.ConfigPath
	subsTest := config.SubscriptionsPath
	listTest := config.ListCredoPath
	pubTest := config.PubCredoPath

	setConfigs()

	if configTest == *configFlag {
		t.Errorf("unexpected result: %s", configTest)
	}

	if subsTest == *subsFlag {
		t.Errorf("unexpected result: %s", subsTest)
	}

	if listTest == *listFlag {
		t.Errorf("unexpected result: %s", listTest)
	}

	if pubTest == *pubFlag {
		t.Errorf("unexpected result: %s", pubTest)
	}
}

func Test_Main(t *testing.T) {
	defer func() {
		os.Unsetenv("SERVICE_UUID")
		os.Unsetenv("SERVICE_HOST")
		os.Unsetenv("NAMESPACE")
		os.Unsetenv("NAMESPACE_LISTENER")
		os.Unsetenv("NAMESPACE_PUBLISHER")
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("SERVICE_PROCESSOR")
		os.Unsetenv("MQTT_LISTENER_URL")
		os.Unsetenv("MQTT_PUBLISHER_URL")
	}()

	wr := new(writer)
	setLog(wr)

	main()
	if !strings.Contains(wr.data, "level=error") {
		t.Errorf("unexpected result: %s", wr.data)
	}
	svr := getMockServer()
	go svr.ListenAndServe(mockURL)
	<-time.After(time.Millisecond * 100)

	setEnv()
	main()
	if !strings.Contains(wr.data, "level=info") {
		t.Errorf("unexpected result: %s", wr.data)
	}
	svr.Close()
}

type writer struct {
	data string
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.data = string(p)
	return len(p), nil
}

func setLog(wr io.Writer) {
	form := new(logrus.TextFormatter)
	logger.Log = logrus.New()
	logger.Log.SetFormatter(form)
	logger.Log.SetOutput(wr)
}

func setEnv() {
	os.Setenv("SERVICE_NAME", "TEST")
	os.Setenv("MQTT_LISTENER_URL", mockURL)
	os.Setenv("MQTT_PUBLISHER_URL", mockURL)
}

const (
	mockURL = "tcp://:15353"
)

// runMockServer creates mock server for testing
func getMockServer() *service.Server {
	return &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
	}
}
