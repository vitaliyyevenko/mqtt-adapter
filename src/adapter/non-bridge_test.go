package adapter

import (
	"testing"
	"os/exec"
	"runtime"
	"strings"
	"os"
	"bytes"
	"bufio"
	"errors"

	"mqtt-adapter/src/logger"

	"github.com/sirupsen/logrus"
)

func TestClient_run(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	loadConf()

	if runtime.GOOS == "windows" {
		t.Skip()
	}

	cmd := exec.Command("ls")
	cl := &client{
		listener:  TestSubscriber{},
		publisher: TestPublisher{},
		command:   cmd,
		topic:     "test_token",
	}
	cl.run()
	if !strings.Contains(wr.data, "Process with PID:") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}

	cmd.Stdout = nil
	cl.run()
	if !strings.Contains(wr.data, "level=error") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
}

func TestClient_runWithEmptyToken(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	loadConf()

	if runtime.GOOS == "windows" {
		t.Skip()
	}

	cmd := exec.Command("ls")
	cl := &client{
		listener:  TestSubscriber{},
		publisher: TestPublisher{},
		command:   cmd,
		topic:     "",
	}
	cl.run()
	if strings.Contains(wr.data, "Process with PID:") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}

	cmd.Stdout = nil
	cl.run()
	if !strings.Contains(wr.data, "level=error") {
		t.Errorf("unexpected result, got: %q", wr.data)
	}
}

func TestClientGetPipes(t *testing.T) {
	logger.Log = &logrus.Logger{}

	testCases := []struct {
		name    string
		needErr bool
		outPipe bool
		errPipe bool
		inPipe  bool
	}{
		{"Test with not error", false, false, false, false},
		{"Test with not nil stdIn", true, false, false, true},
		{"Test with not nil stdErr", true, false, true, false},
		{"Test with not nil stdOut", true, true, false, false},
	}

	for _, tc := range testCases {
		cmd := helperCommand(t, "TEST")
		cl := &client{
			listener:  TestSubscriber{},
			publisher: TestPublisher{},
			command:   cmd,
			topic:     "",
		}

		t.Run(tc.name, func(t *testing.T) {
			if tc.inPipe {
				cmd.Stdin = os.Stdin
			}
			if tc.errPipe {
				cmd.Stderr = os.Stderr
			}
			if tc.outPipe {
				cmd.Stdout = os.Stderr
			}
			o, e, i, err := cl.getPipes()
			if tc.needErr {
				if err == nil || o != nil || e != nil || i != nil {
					t.Error("Expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func helperCommand(t *testing.T, s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestErrorOut(t *testing.T) {
	wr := new(writer)
	setLog(wr)
	var buf bytes.Buffer
	buf.WriteString("hello world")

	scanner := bufio.NewScanner(&buf)
	readStdErr(scanner)
}

func TestLogError(t *testing.T) {
	w := new(writer)
	setLog(w)
	msgErr := "test_error"
	err := errors.New(msgErr)
	logError(0, err)
	if !strings.Contains(w.data, msgErr) {
		t.Errorf("unexpected result: %s", w.data)
	}
}
