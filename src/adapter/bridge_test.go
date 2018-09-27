package adapter

import (
	"testing"
	"mqtt-adapter/src/config"
	"strings"
)

func TestClient_runBridgeWithEmptyTopic(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	loadConf()

	cl := &client{
		listener:  TestSubscriber{needPanic: false},
		publisher: TestPublisher{},
		topic:     "",
	}
	cl.runBridge()
}

func TestClient_runBridgeWithTopic(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	loadConf()

	cl := &client{
		listener:  TestSubscriber{needPanic: false},
		publisher: TestPublisher{},
		topic:     "test_topic",
	}
	cl.runBridge()
}

func TestWithSameNS(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	config.Config = &config.Configuration{
		Same:               true,
		NamespaceListener:  "test",
		NamespacePublisher: "test",
	}
	cl := client{
		listener:  TestSubscriber{},
		publisher: TestPublisher{},
		topic: "test_topic",
	}
	cl.runBridge()
	if !strings.Contains(wr.data, "Listener and Publisher are the same") {
		t.Errorf("unexpected result, got: %s", wr.data)
	}
}
