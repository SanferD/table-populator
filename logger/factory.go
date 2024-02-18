package logger

import (
	"fmt"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
	"github.com/SanferD/table-populator/ioutil"
)

type LoggerKind int

const (
	Error LoggerKind = -1
	Multi LoggerKind = iota
)

func extractLoggerKind(kind string) (LoggerKind, error) {
	switch kind {
	case "multi":
		return Multi, nil
	default:
		return Error, fmt.Errorf("unrecognized logger kind string '%s'", kind)
	}
}

func New(config config.Config) (domain.Logger, error) {
	loggerKind, err := extractLoggerKind(config.LoggerKind)
	switch loggerKind {
	case Multi:
		fo := new(ioutil.StdFileOps)
		lc := new(ioutil.StdLogCreator)
		return InitializeMultiLogger(config, fo, lc)
	case Error:
		fallthrough
	default:
		return nil, fmt.Errorf("error extracting logger kind: %s", err)
	}
}
