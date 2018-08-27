package main

import (
	"log"
	"flag"

	"mqtt-adapter/src/config"
)

func main() {
	configFile := flag.String("config", "package.json", "Configuration file in JSON-format")
	flag.Parse()
	if len(*configFile) > 0 {
		config.ConfigPath = *configFile
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

}
