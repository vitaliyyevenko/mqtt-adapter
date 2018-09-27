package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"mqtt-adapter/src/config"
	"mqtt-adapter/src/mqtt"
	"mqtt-adapter/src/logger"
)

// runBridge starts program in Bridge mode
func (c *client) runBridge() {
	defer c.close()
	if c.topic == "" {
		logger.Log.Error("Microservice cannot start: topic hasn't been initialized")
		return
	}
	if config.Config.Same && config.Config.NamespacePublisher == config.Config.NamespaceListener {
		logger.Log.Error("cannot start microservice: Listener and Publisher are the same")
		return
	}

	topic := fmt.Sprintf("%s/%s", os.Getenv("NAMESPACE_LISTENER"), c.topic)
	msgChan := make(chan string)
	go c.subscribeBridge(msgChan, topic)

	for msg := range msgChan {
		top, err := c.changeTopic(msg)
		if err != nil {
			continue
		}
		c.publisher.Publish(top)
	}
}

// subscribeBridge listen MQTT server
func (c *client) subscribeBridge(msgChan chan string, topic string) {
	c.listener.SubscribeBridge(topic, msgChan)
}

// changeTopic changes topic of MQTT message
func (c *client) changeTopic(msg string) (string, error) {
	message := new(mqtt.Message)
	err := json.Unmarshal([]byte(msg), message)
	if err != nil {
		logger.Log.Warnf("Cannot unmarshal JSON message from Publisher: %q", msg)
		return "", err
	}
	topic := strings.Replace(message.Topic, config.Config.NamespaceListener, config.Config.NamespacePublisher, 1)
	msg = strings.Replace(msg, message.Topic, topic, 1)
	return msg, nil
}

// close disconnects from MQTT server
func (c *client) close() {
	c.listener.Disconnect()
	c.publisher.Disconnect()
}
