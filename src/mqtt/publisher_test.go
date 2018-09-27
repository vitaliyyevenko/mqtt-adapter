package mqtt

import (
	"strings"
	"testing"

	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
)

func TestPublisher_Publish(t *testing.T) {
	logger.Log = &logrus.Logger{}
	testClient := new(TestMQTTClient)
	pub := &publisher{client: testClient}
	testCases := []struct {
		name    string
		needErr bool
		msg     string
	}{
		{"Test with bad json message", true, "{"},
		{"Test with bad json message", false, "{}"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testClient.needErr = tc.needErr
			err := pub.Publish(tc.msg)
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

func TestPublisher_Disconnect(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	testClient := new(TestMQTTClient)
	pub := &publisher{client: testClient}
	pub.Disconnect()
	if !strings.Contains(wr.data, "MQTT Publisher disconnects from server") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
}
