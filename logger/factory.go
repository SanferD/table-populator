package logger

import (
	"fmt"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
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
		return Error, fmt.Errorf("unrecognized logger kind string: %s", kind)
	}
}

func New(config config.Config) (domain.Logger, error) {
	loggerKind, err := extractLoggerKind(config.LoggerKind)
	if err != nil {
		return nil, fmt.Errorf("error extracting logger kind: %s", config.LoggerKind)
	}
	switch loggerKind {
	case Multi:
		return InitializeMultiLogger(config)
	default:
		return nil, fmt.Errorf("unrecognized log kind '%s'", config.LoggerKind)
	}
}
