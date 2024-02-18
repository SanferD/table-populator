package dataio

import (
	"encoding/csv"
	"fmt"

	"github.com/SanferD/table-populator/domain"
	"github.com/SanferD/table-populator/ioutil"
)

type CSVDataIO struct {
	CSVReader ioutil.CSVReader
	CSVWriter ioutil.CSVWriter
}

func InitializeCSVDataIo(fo ioutil.Ops, inputPath, outputPath string) (*CSVDataIO, error) {
	// initialize csv reader
	fdIn, err := fo.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("error opening input csv file: %s", err)
	}

	// initialize csv writer
	fdOut, err := fo.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("error creating output csv file: %s", err)
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
		return fmt.Errorf("error writing record to file: %s", err)
	}
	csvWriter.Flush()
	return nil
}
