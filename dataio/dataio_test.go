package dataio

import (
	"errors"
	"os"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/stretchr/testify/assert"
)

type testDataIOKind struct {
	str  string
	kind DataIOKind
	err  error
}

type testDataIONew struct {
	inputFp  string
	outputFp string
	str      string
	expected interface{}
	err      error
}

func TestExtractDataIOKind(t *testing.T) {
	testDataIOKinds := []testDataIOKind{
		{str: "csv", kind: CSV, err: nil},
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized dataio kind string 'unrecognized'")},
	}

	for _, tc := range testDataIOKinds {
		kind, err := extractDataIOKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

func TestNew(t *testing.T) {
	testCases := []testDataIONew{
		{str: "csv", expected: &CSVDataIO{}, err: nil, inputFp: "csv.csv", outputFp: "out.csv"},
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
