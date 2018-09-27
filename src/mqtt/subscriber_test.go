package mqtt

import (
	"bytes"
	"strings"
	"testing"
	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
)

func TestSubscriber_Subscribe(t *testing.T) {
	logger.Log = &logrus.Logger{}
	testClient := new(TestMQTTClient)
	sub := &subscriber{client: testClient}
	var buf *bytes.Buffer
	testCases := []struct {
		name    string
		needErr bool
	}{
		{"Test with error token", true},
		{"Test with good token", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testClient.needErr = tc.needErr
			sub.Subscribe("", buf)
		})
	}
}

func TestSubscriber_Disconnect(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	testClient := new(TestMQTTClient)
	sub := &subscriber{client: testClient}
	sub.Disconnect()
	if !strings.Contains(wr.data, "MQTT Listener disconnects from server") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
}

func TestHandlers(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	p := make([]byte, 0)
	buf := bytes.NewBuffer(p)
	msgChan := make(chan string, 1)
	defer close(msgChan)
	client := new(TestMQTTClient)
	handler := subsHandler(buf)
	handler(client, TestMessage{})
	if strings.Contains(wr.data, "Received message from MQTT server") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
	handler = subsBridgeHandler(msgChan)
	handler(client, TestMessage{})
	if strings.Contains(wr.data, "Received message from MQTT server") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
}

func TestSubscriber_SubscribeBridge(t *testing.T) {
	logger.Log = &logrus.Logger{}
	testClient := new(TestMQTTClient)
	sub := &subscriber{client: testClient}
	msgChan := make(chan string)
	defer close(msgChan)
	testCases := []struct {
		name    string
		needErr bool
	}{
		{"Test with error token", true},
		{"Test with good token", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testClient.needErr = tc.needErr
			sub.SubscribeBridge("", msgChan)
		})
	}
}
