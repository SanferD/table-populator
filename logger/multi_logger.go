package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/SanferD/table-populator/config"
)

type MultiLogger struct {
	Loggers []*log.Logger
	File    *os.File
}

func InitializeMultiLogger(config config.Config) (*MultiLogger, error) {
	loggers := make([]*log.Logger, 0)

	// initialize stdout logger
	if config.LogToStdout {
		stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)
		loggers = append(loggers, stdoutLogger)
	}

	//initialize file logger
	var file os.File
	if config.LogFilePath != nil {
		file, err := os.Create(*config.LogFilePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file for file logger: %s", err)
		}
		fileLogger := log.New(file, "", log.LstdFlags)
		loggers = append(loggers, fileLogger)
	}
	return &MultiLogger{Loggers: loggers, File: &file}, nil
}

func (multiLogger *MultiLogger) Close() {
	multiLogger.File.Close()
}

func (multiLogger *MultiLogger) Info(msgs ...any) {
	multiLogger.doLog("info", msgs...)
}

func (multiLogger *MultiLogger) Warn(msgs ...any) {
	multiLogger.doLog("warn", msgs...)
}

func (multiLogger *MultiLogger) Debug(msgs ...any) {
	multiLogger.doLog("debug", msgs...)
}

func (multiLogger *MultiLogger) Error(msgs ...any) {
	multiLogger.doLog("error", msgs...)
}

func (multiLogger *MultiLogger) Fatal(msgs ...any) {
	multiLogger.doLog("fatal", msgs...)
}

func (multiLogger *MultiLogger) doLog(prefix string, msgs ...any) {
	for _, logger := range multiLogger.Loggers {
		logger.SetPrefix(prefix + ":")
		logger.Println(msgs...)
	}
	if prefix == "fatal" {
		os.Exit(1)
	}
}
