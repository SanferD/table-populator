package config

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filePath = ".env" // ok, local director, so shouldn't overwrite toplevel config, not thread safe
const (
	csvDataFilePath    = "csv-data-file-path"
	mapsAPIKey         = "maps-api-key"
	outputFilePath     = "output-file-path"
	dataioKind         = "dataio-kind"
	defaultDataIOKind  = "csv"
	locatorKind        = "locator-kind"
	defaultLocatorKind = "google-maps"
	loggerKind         = "logger-kind"
	defaultLoggerKind  = "multi"
	logToStdout        = false
	defaultLogToStdout = true
	logFilePath        = "log-file-path"
)

var defaultLogFilePath *string = nil

const envMinimal = "\n" +
	"CSV_DATA_FILE_PATH=" + csvDataFilePath + "\n" +
	"MAPS_API_KEY=" + mapsAPIKey + "\n" +
	"OUTPUT_FILE_PATH=" + outputFilePath + "\n"

var envFull = "\n" +
	"CSV_DATA_FILE_PATH=" + csvDataFilePath + "\n" +
	"MAPS_API_KEY=" + mapsAPIKey + "\n" +
	"OUTPUT_FILE_PATH=" + outputFilePath + "\n" +
	"DATAIO_KIND=" + dataioKind + "\n" +
	"LOCATOR_KIND=" + locatorKind + "\n" +
	"LOGGER_KIND=" + loggerKind + "\n" +
	"LOG_TO_STDOUT=" + strconv.FormatBool(logToStdout) + "\n" +
	"LOG_FILE_PATH=" + logFilePath + "\n"

const envUnmarshalFailure = "\n" +
	"CSV_DATA_FILE_PTH=" + csvDataFilePath + "\n" +
	"MAPS_API_KEY=" + mapsAPIKey + "\n" +
	"LOG_TO_STDOUT=yes\n"

func TestNew(t *testing.T) {
	t.Run("env full", func(t *testing.T) {
		if err := os.WriteFile(filePath, []byte(envFull), 0644); err != nil {
			assert.NoError(t, err)
		}
		defer os.Remove(filePath)

		config, err := New()

		assert.NoError(t, err)
		assert.Equal(t, csvDataFilePath, config.CSVDataFilePath)
		assert.Equal(t, mapsAPIKey, config.MapsAPIKey)
		assert.Equal(t, outputFilePath, config.OutputCSVFilePath)
		assert.Equal(t, dataioKind, config.DataIOKind)
		assert.Equal(t, locatorKind, config.LocatorKind)
		assert.Equal(t, loggerKind, config.LoggerKind)
		assert.Equal(t, logToStdout, config.LogToStdout)
		assert.Equal(t, logFilePath, *config.LogFilePath)
	})

	t.Run("env minimal", func(t *testing.T) {
		if err := os.WriteFile(filePath, []byte(envMinimal), 0644); err != nil {
			assert.NoError(t, err)
		}
		defer os.Remove(filePath)

		config, err := New()

		assert.NoError(t, err)
		assert.Equal(t, csvDataFilePath, config.CSVDataFilePath)
		assert.Equal(t, mapsAPIKey, config.MapsAPIKey)
		assert.Equal(t, outputFilePath, config.OutputCSVFilePath)
		assert.Equal(t, defaultDataIOKind, config.DataIOKind)
		assert.Equal(t, defaultLocatorKind, config.LocatorKind)
		assert.Equal(t, defaultLoggerKind, config.LoggerKind)
		assert.Equal(t, defaultLogToStdout, config.LogToStdout)
		assert.Nil(t, config.LogFilePath)
	})

	t.Run("env readin failure", func(t *testing.T) {
		_, err := New()
		assert.NotNil(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "error reading in configuration: "))
	})

	t.Run("env unmarshalling failure", func(t *testing.T) {
		if err := os.WriteFile(filePath, []byte(envUnmarshalFailure), 0644); err != nil {
			assert.NoError(t, err)
		}
		defer os.Remove(filePath)

		_, err := New()
		assert.NotNil(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "error unmarshalling configuration"))
	})
}
