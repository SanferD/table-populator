package dataio

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/SanferD/table-populator/domain"
)

type CSVDataIO struct {
	CSVReader *csv.Reader
	CSVWriter *csv.Writer
}

func InitializeCsvDataIo(inputPath string, outputPath string) (*CSVDataIO, error) {
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

	csvDataIO := CSVDataIO{
		CSVReader: csv.NewReader(fdIn),
		CSVWriter: csv.NewWriter(fdOut),
	}
	return &csvDataIO, nil
}

func (csvDataIO *CSVDataIO) ReadRecords() ([]domain.DataRecord, error) {
	csvReader := csvDataIO.CSVReader
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

func (csvDataIO *CSVDataIO) WritePlaceWithCity(placeName string, stateCity domain.StateCity) error {
	csvWriter := csvDataIO.CSVWriter
	row := []string{placeName, stateCity.City, stateCity.State}
	if err := csvWriter.Write(row); err != nil {
		return fmt.Errorf("error witing record to file: %s", err)
	}
	csvWriter.Flush()
	return nil
}
