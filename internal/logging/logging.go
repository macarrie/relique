package logging

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

const TEST_LOG_PATH = "/tmp/relique_tests.log"

func Setup(debug bool, filePath string) {
	logger = log.New()

	var writer io.Writer
	if filePath == "" {
		writer = os.Stdout
	} else {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writer = io.MultiWriter(os.Stdout, file)
		} else {
			writer = os.Stdout
			logger.WithFields(log.Fields{
				"err":  err,
				"path": filePath,
			}).Error("Failed to open log to file, using default stdout")
			fmt.Printf("Failed to open log file '%s', using default stdout: %s", filePath, err)
		}
	}

	logger.Out = writer
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	}
	if debug {
		logger.SetLevel(log.DebugLevel)
		formatter.ForceColors = true
	}

	logger.SetFormatter(formatter)
}

type Fields = log.Fields
type Entry = log.Entry

func WithFields(fields Fields) *log.Entry {
	return logger.WithFields(fields)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}
