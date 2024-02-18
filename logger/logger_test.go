package logger

import (
	"errors"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/stretchr/testify/assert"
)

type testLoggerKind struct {
	str  string
	kind LoggerKind
	err  error
}

type testLoggerNew struct {
	str      string
	expected interface{}
	err      error
}

func TestExtractLoggerKind(t *testing.T) {
	testLoggerKinds := []testLoggerKind{
		{str: "multi", kind: Multi, err: nil},
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized logger kind string 'unrecognized'")},
	}

	for _, tc := range testLoggerKinds {
		kind, err := extractLoggerKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

func TestNewLogger(t *testing.T) {
	testCases := []testLoggerNew{
		{str: "multi", expected: &MultiLogger{}, err: nil},
		{str: "unrecognized", expected: nil, err: errors.New("error extracting logger kind: unrecognized logger kind string 'unrecognized'")},
	}

	for _, tc := range testCases {
		logger, err := New(config.Config{LoggerKind: tc.str})
		assert.IsType(t, tc.expected, logger)
		assert.Equal(t, tc.err, err)
	}
}
