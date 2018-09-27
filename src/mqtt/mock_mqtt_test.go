package mqtt

import (
	"errors"
	"io"
	"time"

	"mqtt-adapter/src/logger"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/surgemq/surgemq/service"
)

type TestMQTTClient struct {
	needErr bool
}

func (t *TestMQTTClient) IsConnected() bool {
	return true
}
func (t *TestMQTTClient) Connect() mqtt.Token {
	return nil
}
func (t *TestMQTTClient) Disconnect(quiesce uint) {}

func (t *TestMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return TestToken{needErr: t.needErr}
}

func (t *TestMQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return TestToken{needErr: t.needErr}
}

func (t *TestMQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return nil
}

func (t *TestMQTTClient) Unsubscribe(topics ...string) mqtt.Token {
	return nil
}
func (t *TestMQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {}

func (t *TestMQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

type TestToken struct {
	needErr bool
}

func (tt TestToken) Wait() bool {
	return true
}
func (tt TestToken) WaitTimeout(time.Duration) bool {
	return false
}
func (tt TestToken) Error() error {
	if tt.needErr {
		return errors.New("test Error")
	}
	return nil
}

type TestMessage struct{}

func (m TestMessage) Duplicate() bool   { return false }
func (m TestMessage) Qos() byte         { return 0 }
func (m TestMessage) Retained() bool    { return false }
func (m TestMessage) Topic() string     { return "" }
func (m TestMessage) MessageID() uint16 { return 0 }
func (m TestMessage) Payload() []byte   { return []byte("test") }

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

const mockURL = "tcp://:15351"

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
