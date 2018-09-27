package adapter

import (
	"bufio"
	"io"
	"fmt"
	"math"
	"syscall"

	"mqtt-adapter/src/config"
	"mqtt-adapter/src/logger"

)

// run launches non-bridge mode
func (c *client) run() {
	outPipe, errPipe, inPipe, err := c.getPipes()
	if err != nil {
		return
	}

	defer func() {
		outPipe.Close()
		errPipe.Close()
		inPipe.Close()
		c.close()
	}()

	scanner := bufio.NewScanner(outPipe)
	scanner.Buffer([]byte{}, math.MaxInt32)
	scannerErr := bufio.NewScanner(errPipe)
	scannerErr.Buffer([]byte{}, math.MaxInt32)
	logger.Log.Infof("Spawning processor: %s", config.Config.ServiceProcessor)
	if err = c.command.Start(); err != nil {
		logger.Log.Error(err)
		return
	}

	pid := c.command.Process.Pid
	logger.Log.Infof("Process with PID: %d has been started", pid)

	defer func() {
		syscall.Kill(pid, 1)
	}()

	if c.topic != "" {
		topic := fmt.Sprintf("%s/%s", config.Config.NamespaceListener, c.topic)
		go c.subscribe(inPipe, topic)
	} else {
		logger.Log.Error("Cannot start Listener: topic is not initialized")
	}

	// read stdOut of the Processor
	go func() {
		for scanner.Scan() {
			go c.publish(scanner.Text())
		}
	}()

	// read stdErr of the Processor
	go readStdErr(scannerErr)

	// wait for Process closes
	err = c.command.Wait()
	logError(c.command.Process.Pid, err)
}

func logError(pid int, err error) {
	if err != nil {
		logger.Log.Warnf("Process with PID (%d) exits: %v", pid, err)
	}
}

// getPipes returns stdOut, stdErr and stdIn of executed process
func (c *client) getPipes() (outPipe, errPipe io.ReadCloser, inPipe io.WriteCloser, err error) {
	outPipe, err = c.command.StdoutPipe()
	if err != nil {
		logger.Log.Errorf("Error obtaining StdOut: %s", err.Error())
		return nil, nil, nil, err
	}
	errPipe, err = c.command.StderrPipe()
	if err != nil {
		logger.Log.Errorf("Error obtaining StdErr: %s", err.Error())
		return nil, nil, nil, err
	}
	inPipe, err = c.command.StdinPipe()
	if err != nil {
		logger.Log.Errorf("Error obtaining StdIn: %s", err.Error())
		return nil, nil, nil, err
	}
	return
}

// subscribe listens to MQTT server
func (c *client) subscribe(w io.Writer, topic string) {
	// for {
	c.listener.Subscribe(topic, w)
	// }
}

func (c *client) publish(msg string) {
	logger.Log.Debugf("processor_stdout_message: %s", msg)
	c.publisher.Publish(msg)
}

func readStdErr(scanner *bufio.Scanner) {
	for scanner.Scan() {
		go logger.Log.Errorf("Processor ERROR event emitted: %s", scanner.Text())
	}
}
