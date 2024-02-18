package locator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"googlemaps.github.io/maps"
)

type testLocatorKind struct {
	str  string
	kind LocatorKind
	err  error
}

func TestExtractLocatorKind(t *testing.T) {
	testLocatorKinds := []testLocatorKind{
		// GoogleMaps kind
		{str: "google-maps", kind: GoogleMaps, err: nil},
		// Error kind
		{str: "unrecognized", kind: Error, err: errors.New("unrecognized locator kind string 'unrecognized'")},
	}

	for _, tc := range testLocatorKinds {
		kind, err := extractLocatorKind(tc.str)
		assert.Equal(t, tc.kind, kind)
		assert.Equal(t, tc.err, err)
	}
}

type testNewLocator struct {
	kind     string
	expected interface{}
	err      error
}

func TestNew(t *testing.T) {
	testNewLocators := []testNewLocator{
		// GoogleMapsLocator
		{kind: "google-maps", expected: &GoogleMapsLocator{}, err: nil},
		// error
		{kind: "unrecognized", expected: nil, err: errors.New("error extracting locator kind: unrecognized locator kind string 'unrecognized'")},
	}
	for _, tc := range testNewLocators {
		config := config.Config{LocatorKind: tc.kind, MapsAPIKey: "api-key"}
		locator, err := New(config)
		assert.IsType(t, tc.expected, locator)
		assert.Equal(t, tc.err, err)
	}
}

type testCaseInitializeMapLocator struct {
	mapsAPIKey  string
	errExpected error
}

func TestInitializeMapLocator(t *testing.T) {
	testCases := []testCaseInitializeMapLocator{
		// valid maps api key
		{mapsAPIKey: "good key", errExpected: nil},
		// invalid maps api key
		{mapsAPIKey: "", errExpected: errors.New("failed to create new maps client: maps: API Key or Maps for Work credentials missing")},
	}

	for _, tc := range testCases {
		config := config.Config{MapsAPIKey: tc.mapsAPIKey}
		_, err := InitializeMapLocator(config)
		assert.Equal(t, tc.errExpected, err)
	}
}

type testCaseGetLocation struct {
	placeName           string
	errTextSearch       error
	placeSearchResponse maps.PlacesSearchResponse
	errExpected         error
	stateCity           *domain.StateCity
}

const badAddress = "bad-address"
const state = "state"
const city = "city"

var badPlaceSearchResponse = maps.PlacesSearchResponse{
	Results: []maps.PlacesSearchResult{
		{FormattedAddress: badAddress},
	},
}

var goodPlaceSearchResponse = maps.PlacesSearchResponse{
	Results: []maps.PlacesSearchResult{
		{FormattedAddress: fmt.Sprintf("1212,%s,%s zip,usa", city, state)},
	},
}

var emptyPlacesSearchResponse = maps.PlacesSearchResponse{
	Results: []maps.PlacesSearchResult{},
}

func TestGetLocation(t *testing.T) {
	testCases := []testCaseGetLocation{
		// TextSearch error
		{placeName: "abc", stateCity: nil, errTextSearch: errors.New("text search error"), placeSearchResponse: emptyPlacesSearchResponse, errExpected: errors.New("failed to get text response for 'abc': text search error")},
		// TextSearch no results
		{placeName: "abc", stateCity: nil, errTextSearch: nil, placeSearchResponse: emptyPlacesSearchResponse, errExpected: errors.New("failed to get any results for 'abc'")},
		// TextSearch malformatted result
		{placeName: "abc", stateCity: nil, errTextSearch: nil, placeSearchResponse: badPlaceSearchResponse, errExpected: fmt.Errorf("formatted address '%s' doesn't have 3 commas", badAddress)},
		// ok
		{placeName: "abc", stateCity: &domain.StateCity{State: state, City: city}, errTextSearch: nil, placeSearchResponse: goodPlaceSearchResponse, errExpected: nil},
	}

	for _, tc := range testCases {
		mockMapsClient := new(MockMapsClient)
		locator := &GoogleMapsLocator{mapsClient: mockMapsClient}
		r := &maps.TextSearchRequest{Query: tc.placeName}
		mockMapsClient.On("TextSearch", mock.Anything, r).Return(tc.placeSearchResponse, tc.errTextSearch)

		stateCity, err := locator.GetLocation(tc.placeName)

		assert.Equal(t, tc.errExpected, err)
		assert.Equal(t, tc.stateCity, stateCity)
	}
}
