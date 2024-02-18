package dataio

import (
	"errors"
	"os"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
	"github.com/SanferD/table-populator/ioutil"
	"github.com/stretchr/testify/assert"
)

var emptyParsedCSVRows [][]string
var parsedCSVRows = [][]string{
	{"ignore-1A", "ignore-1B", "place-1", "ignore-1C"},
	{"ignore-2A", "ignore-2B", "place-2", "ignore-2C"},
	{"ignore-3A", "ignore-3B", "place-3", "ignore-3C"},
}

var emptyDataRecords []domain.DataRecord = []domain.DataRecord{}
var dataRecords = []domain.DataRecord{
	{PlaceName: "place-1"},
	{PlaceName: "place-2"},
	{PlaceName: "place-3"},
}

const (
	csvInputPath  = "input-path"
	csvOutputPath = "output-path"
)

var emptyFile *os.File

type testDataIOKind struct {
	str  string
	kind DataIOKind
	err  error
}

func TestExtractDataIOKind(t *testing.T) {
	testDataIOKinds := []testDataIOKind{
		// csv
		{str: "csv", kind: CSV, err: nil},
		// error
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized dataio kind string 'unrecognized'")},
	}

	for _, tc := range testDataIOKinds {
		kind, err := extractDataIOKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

type testDataIONew struct {
	inputFp  string
	outputFp string
	str      string
	expected interface{}
	err      error
}

func TestNew(t *testing.T) {
	testCases := []testDataIONew{
		// csv
		{str: "csv", expected: &CSVDataIO{}, err: nil, inputFp: "csv.csv", outputFp: "out.csv"},
		// error
		{str: "unrecognized", expected: nil, err: errors.New("error extracting dataio kind: unrecognized dataio kind string 'unrecognized'")},
	}
	for _, tc := range testCases {
		fIn, err := os.CreateTemp("", tc.inputFp)
		assert.Nil(t, err)
		defer fIn.Close()
		fOut, err := os.CreateTemp("", tc.inputFp)
		assert.Nil(t, err)
		defer fOut.Close()
		config := config.Config{DataIOKind: tc.str, CSVDataFilePath: fIn.Name(), OutputCSVFilePath: fOut.Name()}

		dataIO, err := New(config)

		assert.IsType(t, tc.expected, dataIO)
		assert.Equal(t, tc.err, err)
	}
}

type testInitializeCSVDataIO struct {
	createFile  *os.File
	openFile    *os.File
	createErr   error
	openErr     error
	expectedErr error
}

func TestInitializeCSVDataIO(t *testing.T) {
	testCases := []testInitializeCSVDataIO{
		// OK
		{createFile: emptyFile, openFile: emptyFile, createErr: nil, openErr: nil, expectedErr: nil},
		// Open error
		{createFile: nil, openFile: nil, createErr: nil, openErr: errors.New("open error"), expectedErr: errors.New("error opening input csv file: open error")},
		// Create error
		{createFile: emptyFile, openFile: nil, createErr: errors.New("create error"), openErr: nil, expectedErr: errors.New("error creating output csv file: create error")},
	}
	for _, tc := range testCases {
		mockIOUtil := new(ioutil.MockIOUtil)
		mockIOUtil.On("Open", csvInputPath).Return(tc.openFile, tc.openErr)
		mockIOUtil.On("Create", csvOutputPath).Return(tc.createFile, tc.createErr)

		_, err := InitializeCSVDataIo(mockIOUtil, csvInputPath, csvOutputPath)

		assert.Equal(t, tc.expectedErr, err)
	}
}

type testReadRecords struct {
	errReadAll  error
	retReadAll  [][]string
	retExpected []domain.DataRecord
	errExpected error
}

func TestReadRecrods(t *testing.T) {
	testCases := []testReadRecords{
		// OK
		{errReadAll: nil, retReadAll: parsedCSVRows, retExpected: dataRecords, errExpected: nil},
		// ReadAll error
		{errReadAll: errors.New("read all error"), retReadAll: emptyParsedCSVRows, retExpected: emptyDataRecords, errExpected: errors.New("error with csv reader readall: read all error")},
	}

	for _, tc := range testCases {
		mockCSVReader := new(ioutil.MockCSVReader)
		mockCSVWriter := new(ioutil.MockCSVWriter)
		mockCSVReader.On("ReadAll", nil).Return(tc.retReadAll, tc.errReadAll)

		csvDataIO := CSVDataIO{CSVReader: mockCSVReader, CSVWriter: mockCSVWriter}
		dataRecords, err := csvDataIO.ReadRecords()

		assert.Equal(t, tc.retExpected, dataRecords)
		assert.Equal(t, tc.errExpected, err)
	}
}

type testWriteRecords struct {
	placeName               string
	stateCity               domain.StateCity
	errWrite                error
	errExpected             error
	flushExpectedCallsCount int
}

func TestWriteRecords(t *testing.T) {
	testCases := []testWriteRecords{
		// OK
		{placeName: "place-name", stateCity: domain.StateCity{State: "state", City: "city"}, errWrite: errors.New("write error"), errExpected: errors.New("error writing record to file: write error"), flushExpectedCallsCount: 0},
		// Write error
		{placeName: "place-name", stateCity: domain.StateCity{State: "state", City: "city"}, errWrite: nil, errExpected: nil, flushExpectedCallsCount: 1},
	}

	for _, tc := range testCases {
		mockCSVReader := new(ioutil.MockCSVReader)
		mockCSVWriter := new(ioutil.MockCSVWriter)
		row := []string{tc.placeName, tc.stateCity.City, tc.stateCity.State}
		mockCSVWriter.On("Write", row).Return(tc.errWrite)
		mockCSVWriter.On("Flush", nil).Return(nil)

		csvDataIO := CSVDataIO{CSVReader: mockCSVReader, CSVWriter: mockCSVWriter}
		err := csvDataIO.WritePlaceWithCity(tc.placeName, tc.stateCity)

		assert.Equal(t, tc.errExpected, err)
		mockCSVWriter.AssertNumberOfCalls(t, "Flush", tc.flushExpectedCallsCount)
		mockCSVWriter.AssertNumberOfCalls(t, "Write", 1)
	}
}
