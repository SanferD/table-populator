package dataio

import (
	"fmt"
	"github.com/SanferD/table-populator/domain"
	configMod "github.com/SanferD/table-populator/config"
)

type DataIoKind int
const (
	Error 	DataIoKind = -1
	Csv 	DataIoKind = iota	
)

func extractDataIoKind(kind string) (DataIoKind, error) {
	switch kind {
	case "csv":
		return Csv, nil
	default:
		return Error, fmt.Errorf("unrecognized dataio kind string: %s", kind)
	}
}

func CreateDataIo(config configMod.Config) (domain.DataIo, error) {
	dataIoKind, err := extractDataIoKind(config.DataIoKind)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataio kind: %s", err)
	}
	switch dataIoKind {
	case Csv:
		return InitializeCsvDataIo(config.CsvDataFilePath, config.OutputCsvFilePath)
	default:
		return nil, fmt.Errorf("unrecognized kind '%s'", config.DataIoKind)
	}
}
