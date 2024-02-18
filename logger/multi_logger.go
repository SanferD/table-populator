package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/ioutil"
)

type MultiLogger struct {
	loggers []ioutil.Logger
	file    io.Closer
	fileOps ioutil.FileOps
}

func InitializeMultiLogger(config config.Config, fo ioutil.FileOps, lc ioutil.LogCreator) (*MultiLogger, error) {
	loggers := make([]ioutil.Logger, 0)

	// initialize stdout logger
	if config.LogToStdout {
		stdoutLogger := lc.New(os.Stdout, "", log.LstdFlags)
		loggers = append(loggers, stdoutLogger)
	}

	// initialize file logger
	var file os.File
	if config.LogFilePath != nil {
		file, err := fo.Create(*config.LogFilePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file for file logger: %s", err)
		}
		fileLogger := lc.New(file, "", log.LstdFlags)
		loggers = append(loggers, fileLogger)
	}
	return &MultiLogger{loggers: loggers, file: &file, fileOps: fo}, nil
}

func (multiLogger *MultiLogger) Close() error {
	if err := multiLogger.file.Close(); err != nil {
		return fmt.Errorf("error closing file: %s", err)
	}
	return nil
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
	for _, logger := range multiLogger.loggers {
		logger.SetPrefix(prefix + ":")
		logger.Println(msgs...)
	}
	if prefix == "fatal" {
		multiLogger.fileOps.Exit(1)
	}
}
