package dataio

import (
	"fmt"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
)

type DataIOKind int

const (
	Error DataIOKind = -1
	CSV   DataIOKind = iota
)

func extractDataIOKind(kind string) (DataIOKind, error) {
	switch kind {
	case "csv":
		return CSV, nil
	default:
		return Error, fmt.Errorf("unrecognized dataio kind string '%s'", kind)
	}
}

func New(config config.Config) (domain.DataIO, error) {
	dataIOKind, err := extractDataIOKind(config.DataIOKind)
	switch dataIOKind {
	case CSV:
		return InitializeCsvDataIo(config.CSVDataFilePath, config.OutputCSVFilePath)
	case Error:
		fallthrough
	default:
		return nil, fmt.Errorf("error extracting dataio kind: %s", err)
	}
}
