package mqtt

import (
	"encoding/json"
	"mqtt-adapter/src/logger"

	"github.com/eclipse/paho.mqtt.golang"
)

// publisher is an instance of Publisher interface
type publisher struct {
	client mqtt.Client
}

// Publish publishes specified message to MQTT server
func (p *publisher) Publish(msg string) error {
	m := new(Message)
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		logger.Log.Warnf("Cannot unmarshal JSON message from Process: %q", msg)
		return err
	}
	topic := m.Topic
	token := p.client.Publish(topic, qos, false, msg)
	token.Wait()
	return token.Error()
}

// Disconnect ends the connection with the server
func (p *publisher) Disconnect() {
	if p.client.IsConnected() {
		logger.Log.Infoln("MQTT Publisher disconnects from server")
		p.client.Disconnect(250)
	}
}
