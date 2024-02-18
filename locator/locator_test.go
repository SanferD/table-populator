package locator

import (
	"errors"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/stretchr/testify/assert"
)

type testLocatorKind struct {
	str  string
	kind LocatorKind
	err  error
}

type testNewLocator struct {
	kind     string
	expected interface{}
	err      error
}

func TestExtractLocatorKind(t *testing.T) {
	testLocatorKinds := []testLocatorKind{
		{str: "google-maps", kind: GoogleMaps, err: nil},
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized locator kind string 'unrecognized'")},
	}

	for _, tc := range testLocatorKinds {
		kind, err := extractLocatorKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

func TestNew(t *testing.T) {
	testNewLocators := []testNewLocator{
		{kind: "google-maps", expected: &GoogleMapsLocator{}, err: nil},
		{kind: "unrecognized", expected: nil, err: errors.New("error extracting locator kind: unrecognized locator kind string 'unrecognized'")},
	}
	for _, tc := range testNewLocators {
		config := config.Config{LocatorKind: tc.kind, MapsAPIKey: "api-key"}
		locator, err := New(config)
		assert.IsType(t, tc.expected, locator)
		assert.Equal(t, tc.err, err)
	}
}
