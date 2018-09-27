# mqtt-adapter

Golang implementation of MQTT service adapter

### Dependencies


To install glide follow link [Glide install](https://github.com/Masterminds/glide "Glide install")

To install Gometalinter follow link [Gometalinter.V2 install](https://github.com/alecthomas/gometalinter "Gometalinter.V2 install")

### Use make file

Make sure, that current directory is `src/`

To download dependencies follow next command
```bash
/src $ make dependencies
```

To show test coverage follow next command
```bash
/src $ make test
```

To check code quality follow next command
```bash
/src $ make code-quality
```
There would be file `static-analysis.xml` in root directory

To build binary file folow next command:
```bash
/src $ make build
```

or

```
/src $ make build-mac
```

You will find `microservice-adapter-mqtt` in `../dev/` directory

To run microservice-adapter-mqtt, make sure that all environments were set:
```bash
$SERVICE_NAME
$MQTT_LISTENER_URL
$MQTT_PUBLISHER_URL
$SERVICE_PROCESSOR
$DEBUG
$BRIDGE
$NAMESPACE
$NAMESPACE_LISTENER
$NAMESPACE_PUBLISHER
```
Examples of setting `$SERVICE_PROCESSOR` :
```bash
export SERVICE_PROCESSOR=./service-processor/processor
 ...
 export SERVICE_PROCESSOR="go run processor.go"
  ...
 export SERVICE_PROCESSOR="go run processor.go --conf=../path/to/config/processor"
  ...
 export SERVICE_PROCESSOR="ls -lah"
 ...
 export SERVICE_PROCESSOR="node hello.js"
```

To launch `microservice-adapter-mqtt` follow next command:
```
 microservice-adapter-mqtt --conf=path/to/package.json --subs=path/to/subscriptions.txt --list=path/to/mqtt_listener.json --pub=path/to/mqtt_publisher.json
```

If any flags won't be set, `microservice-adapter-mqtt` uses default paths:
```bash
path/to/package.json = ./package.json
path/to/subscriptions.txt = ./service-processor/subscriptions.txt
path/to/mqtt_listener.json = /run/secrets/mqtt_listener.json
path/to/mqtt_publisher.json = /run/secrets/mqtt_publisher.json
```

To launch example with processor follow next command:

```
make example
```


