package logger

import (
	"github.com/sirupsen/logrus"
)

const (
	// logTimeFormat represents time format in log messages
	logTimeFormat = "2006-01-02 15:04:05.99"
)

// Log is an instance to log messages
var Log *logrus.Logger

func init() {
	initLogger()
}

// initLogger initializes Log with default options
func initLogger(){
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = logTimeFormat
	customFormatter.FullTimestamp = true
	Log = logrus.New()
	Log.SetFormatter(customFormatter)
}
