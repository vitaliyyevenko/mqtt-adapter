package adapter

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"mqtt-adapter/src/config"
	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	svr := getMockServer()
	defer svr.Close()
	go svr.ListenAndServe(mockURL)
	<-time.After(time.Millisecond * 100)

	logger.Log = &logrus.Logger{}
	testCases := []struct {
		name        string
		needErr     bool
		loadConfig  bool
		breakPubURL bool
		breakSubURL bool
	}{
		{"Test with <nil> Config", true, false, false, false},
		{"Test with not <nil> Config", false, true, false, false},
		{"Test with bad Publisher URL", true, true, true, false},
		{"Test with bad Subscriber URL", true, true, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.loadConfig {
				config.Config = nil
			} else {
				loadConf()
			}
			if tc.breakPubURL {
				config.Config.MQTTPublisherURL = ""
			}
			if tc.breakSubURL {
				config.Config.MQTTListenerURL = ""
			}
			_, err := New()
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

func loadConf() {
	config.Config = new(config.Configuration)
	config.Config.MQTTPublisherURL = mockURL
	config.Config.MQTTListenerURL = mockURL
	config.Config.ServiceProcessor = "some_test_process"
}

func TestClient_Run(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	loadConf()
	config.Config.Bridge = false
	cl := &client{
		listener:  TestSubscriber{needPanic: true},
		publisher: TestPublisher{},
		command:   new(exec.Cmd),
	}
	cl.Run()
	if !strings.Contains(wr.data, "level=error") {
		t.Errorf("unexpected result: expected: 'level=error', got: %q", wr.data)
	}
	config.Config.Bridge = true
	cl.Run()
}
