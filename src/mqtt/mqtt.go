package mqtt

import (
	"fmt"
	"encoding/json"

	"github.com/eclipse/paho.mqtt.golang"
)

type mqttClient struct {
	id     string
	topic  string
	client mqtt.Client
}

type message struct {
	Topic   string  `json:"topic,omitempty"`
	Name    string  `json:"service_name,omitempty"`
	UUID    string  `json:"service_uuid,omitempty"`
	Host    string  `json:"service_host,omitempty"`
	Payload payLoad `json:"payload,omitempty"`
}

type payLoad struct {
	LogEntry logEntry `json:"log_entry,omitempty"`
}

type logEntry struct {
	Level   string `json:"log_level,omitempty"`
	Message string `json:"log_message,omitempty"`
}

func (m message) String() string {
	p, _ := json.Marshal(&m)
	return string(p)
}

/*func NewMQTTClient(level, username, pass string) mqtt.Client {
	c := new(mqttClient)
	id := fmt.Sprintf("%s_%s_%s_listener", Config.Name, Config.Host, Config.UUID)
	c.topic = fmt.Sprintf("%s/log/%s/%s/%s", Config.NamespacePublisher, Config.Name, Config.UUID, level)
	options := mqtt.NewClientOptions()
	options.AddBroker(Config.MQTTListenerURL)
	options.SetClientID(id)
	// options.SetWill()

	c.client = mqtt.NewClient(options)
	return c.client
}*/

func main() {
	log := logEntry{
		Level:   "someLevel",
		Message: "someMessage",
	}

	m := message{
		UUID:    "someUUID",
		Name:    "someName",
		Host:    "someHost",
		Topic:   "someTopic",
		Payload: payLoad{log},
	}
	fmt.Println(m)
}
