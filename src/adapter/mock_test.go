package adapter

import (
	"io"

	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
	"github.com/surgemq/surgemq/service"
)

type TestSubscriber struct {
	needPanic bool
}

func (s TestSubscriber) Subscribe(topic string, writer io.Writer) {}

func (s TestSubscriber) SubscribeBridge(topic string, msgChan chan<- string) {
	if s.needPanic {
		panic("test Panic")
	}
	msgChan <- `{"topic":"test"}`
	msgChan <- `{"topic":123}`
	close(msgChan)
}

func (s TestSubscriber) Disconnect() {}

type TestPublisher struct{}

func (p TestPublisher) Publish(msg string) error { return nil }

func (p TestPublisher) Disconnect() {}

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

const mockURL = "tcp://:15352"

// runMockServer creates mock server for testing
func getMockServer()*service.Server {
	return  &service.Server{
		KeepAlive:        300,           // seconds
		ConnectTimeout:   2,             // seconds
		SessionsProvider: "mem",         // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
	}
}