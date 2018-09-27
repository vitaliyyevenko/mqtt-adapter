package mqtt

import (
	"testing"
	"time"
	"mqtt-adapter/src/config"
	"mqtt-adapter/src/logger"
	"github.com/sirupsen/logrus"
)

func TestNewClient(t *testing.T) {
	svr := getMockServer()
	defer svr.Close()
	go svr.ListenAndServe(mockURL)
	<-time.After(time.Millisecond * 100)

	testCases := []struct {
		name      string
		needErr   bool
		brokerURL string
		userName  string
	}{
		{"Test with empty broker", true, "", ""},
		{"Test with empty broker and not empty user", true, "", "test"},
		{"Test with not empty broker", false, mockURL, ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newClient(tc.brokerURL, "test", config.Credentials{UserName: tc.userName})
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

func TestNewMQTTClients(t *testing.T) {
	svr := getMockServer()
	defer svr.Close()
	go svr.ListenAndServe(mockURL)
	<-time.After(time.Millisecond * 100)

	logger.Log = &logrus.Logger{}

	c := &config.Configuration{}

	_, _, err := NewMQTTClients(c)
	if err == nil {
		t.Error("expected not nil error")
	}
	c.MQTTListenerURL = mockURL
	c.Same = true
	_, _, err = NewMQTTClients(c)
	if err != nil {
		t.Error("expected nil error")
	}

	c.Same = false
	_, _, err = NewMQTTClients(c)
	if err == nil {
		t.Error("expected not nil error")
	}

	c.MQTTPublisherURL = mockURL
	_, _, err = NewMQTTClients(c)
	if err != nil {
		t.Error("expected nil error")
	}
}
