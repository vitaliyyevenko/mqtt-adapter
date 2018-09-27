package mqtt

import (
	"io"
	"fmt"
	"time"

	"mqtt-adapter/src/logger"

	"github.com/eclipse/paho.mqtt.golang"
)

// is an instance of Subscriber interface
type subscriber struct {
	client mqtt.Client
}

var (
	subsHandler = func(writer io.Writer) func(client mqtt.Client, msg mqtt.Message) {
		return func(client mqtt.Client, msg mqtt.Message) {
			logger.Log.Debugf("MQTT_MESSAGE_RECEIVED: %s", msg.Payload())
			fmt.Fprintln(writer, string(msg.Payload()))
		}
	}

	subsBridgeHandler = func(msgChan chan<- string) func(client mqtt.Client, msg mqtt.Message) {
		return func(client mqtt.Client, msg mqtt.Message) {
			logger.Log.Debugf("MQTT message relayed through bridge: %s", msg.Payload())
			msgChan <- string(msg.Payload())
		}
	}
)

// Subscribe starts a new subscription in non-bridge mode and writs received message to io.Writer
func (s *subscriber) Subscribe(topic string, writer io.Writer) {

	if token := s.client.Subscribe(topic, qos, subsHandler(writer)); token.Wait() && token.Error() != nil {
		time.Sleep(time.Millisecond * 10)
		return
	}
}

// SubscribeBridge starts a new subscription in non-bridge mode and writs received message to specified channel
func (s *subscriber) SubscribeBridge(topic string, msgChan chan<- string) {
	if token := s.client.Subscribe(topic, qos, subsBridgeHandler(msgChan)); token.Wait() && token.Error() != nil {
		time.Sleep(time.Millisecond * 10)
		return
	}
}

// Disconnect ends the connection with the server
func (s *subscriber) Disconnect() {
	if s.client.IsConnected() {
		logger.Log.Infoln("MQTT Listener disconnects from server")
		s.client.Disconnect(250)
	}
}
