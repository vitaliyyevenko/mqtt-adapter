package mqtt

import (
	"io"
	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
	"mqtt-adapter/src/config"
)

const qos = 0

// Subscriber is an interface that describes behavior of a subscriber to MQTT
type Subscriber interface {
	Subscribe(topic string, writer io.Writer)
	SubscribeBridge(topic string, msgChan chan<- string)
	Disconnect()
}

// Publisher is an interface that describes behavior of a publisher to MQTT
type Publisher interface {
	Publish(msg string) error
	Disconnect()
}

// Message represent model of MQTT message
type Message struct {
	Topic string `json:"topic"`
}

func newClient(broker, clientID string, credo config.Credentials) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetCleanSession(true)
	opts.SetClientID(clientID)
	if credo.UserName != "" || credo.Password != "" {
		opts.SetUsername(credo.UserName)
		opts.SetPassword(credo.Password)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("cannot connect to MQTT broker (%s): %v", broker, token.Error())
	}
	return client, nil
}

// NewMQTTClients creates and initializes publisher and listener
func NewMQTTClients(conf *config.Configuration) (pub Publisher, sub Subscriber, err error) {
	var clS, clP mqtt.Client
	listClientID := fmt.Sprintf("%s_%s_%s_lis", conf.Name, conf.Host, conf.UUID)
	clS, err = newClient(conf.MQTTListenerURL, listClientID, conf.ListCredo)
	if err != nil {
		return nil, nil, err
	}
	if conf.Same {
		sub = &subscriber{client: clS}
		pub = &publisher{client: clS}
		return pub, sub, nil
	}
	pubClientID := fmt.Sprintf("%s_%s_%s_pub", conf.Name, conf.Host, conf.UUID)
	clP, err = newClient(conf.MQTTPublisherURL, pubClientID, conf.PubCredo)
	if err != nil {
		return nil, nil, err
	}
	sub = &subscriber{client: clS}
	pub = &publisher{client: clP}
	return pub, sub, nil
}
