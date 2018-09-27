package main

import (
	"flag"

	"mqtt-adapter/src/adapter"
	"mqtt-adapter/src/config"
	"mqtt-adapter/src/logger"
)

var (
	configFlag = flag.String("conf", "./package.json", "Path to package.json file")
	subsFlag   = flag.String("subs", "./service-processor/subscriptions.txt", "Path to subscriptions.txt file")
	listFlag   = flag.String("list", "/run/secrets/mqtt_listener.json", "Path to mqtt_listener.json file")
	pubFlag    = flag.String("pub", "/run/secrets/mqtt_publisher.json", "Path to mqtt_publisher.json file")
)

func main() {
	setConfigs()

	logger.Log.Infoln("Start MicroService MQTT adapter ...")
	err := config.Load()
	if err != nil {
		logger.Log.Error(err)
		return
	}
	ms, err := adapter.New()
	if err != nil {
		logger.Log.Error(err)
		return
	}

	ms.Run()
}

// setConfigs sets paths to files from command line
func setConfigs() {
	flag.Parse()

	config.ConfigPath = *configFlag
	config.SubscriptionsPath = *subsFlag
	config.ListCredoPath = *listFlag
	config.PubCredoPath = *pubFlag
}
