package dataio

import (
	"os"
	"fmt"
	"encoding/csv"

	"github.com/SanferD/table-populator/domain"
)

type CsvDataIo struct {
	CsvReader *csv.Reader
	CsvWriter *csv.Writer
}

func InitializeCsvDataIo(inputPath string, outputPath string) (*CsvDataIo, error) {
	// initialize csv reader
	fdIn, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("error opening input csv file: %s", err)
	}

	// initialize csv writer
	fdOut, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("error opening output csv file: %s", err)
	}

	csvDataIo := CsvDataIo{
		CsvReader: csv.NewReader(fdIn),
		CsvWriter: csv.NewWriter(fdOut),
	}
	return &csvDataIo, nil
}

func (csvDataIo *CsvDataIo) ReadRecords() ([]domain.DataRecord, error) {
	csvReader := csvDataIo.CsvReader
	var dataRecords = make([]domain.DataRecord, 0)
	records, err := csvReader.ReadAll()
	if err != nil {
		return dataRecords, fmt.Errorf("error with csv reader readall: %s", err)
	}
	for _, record := range records {
		dataRecord := domain.DataRecord{PlaceName: record[2]}
		dataRecords = append(dataRecords, dataRecord)
	}
	return dataRecords, nil
}

func (csvDataIo *CsvDataIo) WritePlaceWithCity(placeName string, stateCity domain.StateCity) error {
	csvWriter := csvDataIo.CsvWriter
	row := []string{placeName, stateCity.City, stateCity.State,}
	if err := csvWriter.Write(row); err != nil {
		return fmt.Errorf("error witing record to file: %s", err)
	}
	csvWriter.Flush()
	return nil
}
