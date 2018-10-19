package log

import (
	"os"
	"fmt"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func EnableDebug() {
	log.SetLevel(logrus.DebugLevel)
}

func init() {
	formatter := &logrus.TextFormatter{
	    FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

func Info(format string, args ...interface{}) {
	log.Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	log.Warn(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(format, args...))
}
