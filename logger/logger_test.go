package logger

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/ioutil"
	"github.com/stretchr/testify/assert"
)

type testLoggerKind struct {
	str  string
	kind LoggerKind
	err  error
}

func TestExtractLoggerKind(t *testing.T) {
	testLoggerKinds := []testLoggerKind{
		// Multi kind
		{str: "multi", kind: Multi, err: nil},
		// Error kind
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized logger kind string 'unrecognized'")},
	}

	for _, tc := range testLoggerKinds {
		kind, err := extractLoggerKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

type testLoggerNew struct {
	str      string
	expected interface{}
	err      error
}

func TestNewLogger(t *testing.T) {
	testCases := []testLoggerNew{
		// Multi logger
		{str: "multi", expected: &MultiLogger{}, err: nil},
		// error
		{str: "unrecognized", expected: nil, err: errors.New("error extracting logger kind: unrecognized logger kind string 'unrecognized'")},
	}

	for _, tc := range testCases {
		logger, err := New(config.Config{LoggerKind: tc.str})
		assert.IsType(t, tc.expected, logger)
		assert.Equal(t, tc.err, err)
	}
}

type testCaseInitializeMultiLogger struct {
	logToStdout      bool
	logFilePath      *string
	errExpected      error
	logCountExpected int // -1 => nil MultiLogger
	errCreate        error
}

var logFilePath = "abc"
var emptyFile *os.File = new(os.File)
var emptyLogger *log.Logger

func TestInitializeMultiLogger(t *testing.T) {
	testCases := []testCaseInitializeMultiLogger{
		// no log file path
		{logToStdout: false, logFilePath: nil, errExpected: nil, logCountExpected: 0, errCreate: nil},
		// log to stdout only
		{logToStdout: true, logFilePath: nil, errExpected: nil, logCountExpected: 1, errCreate: nil},
		// log to file, create error
		{logToStdout: false, logFilePath: &logFilePath, errExpected: errors.New("error creating file for file logger: create error"), logCountExpected: -1, errCreate: errors.New("create error")},
		// log to stdout and file
		{logToStdout: true, logFilePath: &logFilePath, errExpected: nil, logCountExpected: 2, errCreate: nil},
	}

	for _, tc := range testCases {
		mockFileOps := new(ioutil.MockFileOps)
		mockLogCreator := new(ioutil.MockLogCreator)
		if tc.logFilePath != nil {
			mockFileOps.On("Create", *tc.logFilePath).Return(emptyFile, tc.errCreate)
			if tc.errCreate == nil {
				mockLogCreator.On("New", emptyFile, "", log.LstdFlags).Return(emptyLogger)
			}
		}
		if tc.logToStdout {
			mockLogCreator.On("New", os.Stdout, "", log.LstdFlags).Return(emptyLogger)
		}
		config := config.Config{LogToStdout: tc.logToStdout, LogFilePath: tc.logFilePath}

		multiLogger, err := InitializeMultiLogger(config, mockFileOps, mockLogCreator)

		assert.Equal(t, tc.errExpected, err)
		if tc.logCountExpected == -1 {
			assert.Nil(t, multiLogger)
		} else {
			assert.Equal(t, tc.logCountExpected, len(multiLogger.loggers))
		}
		mockFileOps.AssertExpectations(t)
		mockLogCreator.AssertExpectations(t)
	}
}

var msgs = []interface{}{"hello", "world"}

func TestLogging(t *testing.T) {
	prefixes := []string{"info", "warn", "debug", "error", "fatal"}
	mockFileOps := new(ioutil.MockFileOps)
	mockLogger := new(ioutil.MockLogger)
	multiLogger := MultiLogger{
		loggers: []ioutil.Logger{mockLogger},
		fileOps: mockFileOps,
	}
	mockFileOps.On("Exit", 1)
	mockLogger.On("Println", msgs...)
	for _, prefix := range prefixes {
		mockLogger.On("SetPrefix", prefix+":")
	}

	multiLogger.Info(msgs...)
	multiLogger.Warn(msgs...)
	multiLogger.Debug(msgs...)
	multiLogger.Error(msgs...)
	multiLogger.Fatal(msgs...)

	mockLogger.AssertExpectations(t)
}

type testCaseClose struct {
	errClose    error
	errExpected error
}

func TestClose(t *testing.T) {
	testCases := []testCaseClose{
		// Close error
		{errClose: errors.New("close error"), errExpected: errors.New("error closing file: close error")},
		// ok
		{errClose: nil, errExpected: nil},
	}
	for _, tc := range testCases {
		mockCloser := new(ioutil.MockCloser)
		multiLogger := MultiLogger{file: mockCloser}
		mockCloser.On("Close").Return(tc.errClose)

		err := multiLogger.Close()

		mockCloser.AssertExpectations(t)
		assert.Equal(t, tc.errExpected, err)
	}
}
