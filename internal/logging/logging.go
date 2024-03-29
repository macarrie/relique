package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

const TEST_LOG_PATH = "/relique_tests.log"

var IsTest = false

func GetLogRoot() string {
	if IsTest {
		return "/tmp/"
	}

	logRoot := os.Getenv("RELIQUE_LOG_ROOT")
	if logRoot == "" {
		return "/var/log/relique/"
	} else {
		return filepath.Clean(logRoot)
	}
}

func SetupCliLogger(debug bool, outputAsJson bool) {
	logger = log.New()

	logger.Out = os.Stdout
	var formatter log.Formatter
	if outputAsJson {
		formatter = &log.JSONFormatter{}
		// Do not log when is json mode to allow easy parsing by external tools
		logger.Out = io.Discard
	} else {
		formatter = &log.TextFormatter{
			DisableTimestamp: true,
			PadLevelText:     true,
			ForceColors:      true,
		}
	}

	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	logger.SetFormatter(formatter)
}

func Setup(debug bool, logPath string) {
	logger = log.New()

	var writer io.Writer
	if logPath == "" {
		writer = os.Stdout
	} else {
		filePath := filepath.Clean(fmt.Sprintf("%s/%s", GetLogRoot(), logPath))
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
