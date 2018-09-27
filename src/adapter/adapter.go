package adapter

import (
	"fmt"
	"os/exec"
	"strings"

	"mqtt-adapter/src/config"
	"mqtt-adapter/src/mqtt"
	"mqtt-adapter/src/logger"
)

// Runner is a client for Microservice MQTT Adapter
type Runner interface {
	Run()
}

// client is an instance of Microservice MQTT Adapter
type client struct {
	topic     string
	listener  mqtt.Subscriber
	publisher mqtt.Publisher
	command   *exec.Cmd
}

// New initializes MQTT adapter and return instance
func New() (Runner, error) {
	if config.Config == nil {
		return nil, fmt.Errorf("not initialized Config")
	}
	commands := strings.Fields(config.Config.ServiceProcessor)
	adapter := new(client)
	adapter.topic = config.Config.Topic
	pub, sub, err := mqtt.NewMQTTClients(config.Config)
	if err != nil {
		return nil, err
	}
	adapter.publisher = pub
	adapter.listener = sub
	adapter.command = exec.Command(commands[0], commands[1:]...)
	return adapter, nil
}


// Run starts app
func (c *client) Run() {
	defer c.listener.Disconnect()
	if config.Config.Bridge {
		logger.Log.Infoln("Start in Bridge mode")
		c.runBridge()
	} else {
		logger.Log.Infoln("Start in non-Bridge mode")
		c.run()
	}
}
