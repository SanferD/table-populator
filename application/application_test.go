package application

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/SanferD/table-populator/dataio"
	"github.com/SanferD/table-populator/domain"
	"github.com/SanferD/table-populator/locator"
	"github.com/SanferD/table-populator/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var drValues = []domain.DataRecord{
	{PlaceName: "place1"}, {PlaceName: "place2"}, {PlaceName: "place3"},
}

var pnscValues = []PlaceNameStateCity{
	{PlaceName: "place1", StateCity: domain.StateCity{State: "state1", City: "city1"}},
	{PlaceName: "place2", StateCity: domain.StateCity{State: "state2", City: "city2"}},
	{PlaceName: "place3", StateCity: domain.StateCity{State: "state3", City: "city3"}},
}

var errReadRecords = errors.New("read records error")
var errGetLocation = errors.New("locator error on get location")

func TestPopulateOutputWhenNoValues(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	mockDataIO := new(dataio.MockDataIO)
	pnscCh := newPlaceNameStateCityChan()
	dpoCh := newDonePopulateOutputChannel()

	go populateOutput(mockLogger, mockDataIO, pnscCh, dpoCh)
	close(pnscCh)
	done := <-dpoCh
	close(dpoCh)

	assert.True(t, done)
	mockLogger.AssertNotCalled(t, "Debug")
	mockDataIO.AssertNotCalled(t, "WritePlaceWithCity")
}

func TestPopulateOutputWhenValues(t *testing.T) {
	mockLogger := newMockLogger()
	mockDataIO := newMockDataIO(false)
	pnscCh := newPlaceNameStateCityChan()
	dpoCh := newDonePopulateOutputChannel()

	go populateOutput(mockLogger, mockDataIO, pnscCh, dpoCh)
	for _, pnsc := range pnscValues {
		pnscCh <- pnsc
	}
	close(pnscCh)
	done := <-dpoCh
	close(dpoCh)

	assert.True(t, done)
	for _, pnsc := range pnscValues {
		mockLogger.AssertCalled(t, "Debug", "writing place with city", pnsc.PlaceName, pnsc.StateCity)
		mockDataIO.AssertCalled(t, "WritePlaceWithCity", pnsc.PlaceName, pnsc.StateCity)
	}
}

func TestFetchLocations(t *testing.T) {
	mockLogger := newMockLogger()
	mockDataIO := newMockDataIO(false)
	mockLocator := newMockLocator()
	drCh := newDataRecordChan()
	pnscCh := newPlaceNameStateCityChan()
	var wg sync.WaitGroup
	id := 1

	wg.Add(1)
	values := make([]PlaceNameStateCity, 0)
	go func() {
		for pnsc := range pnscCh {
			values = append(values, PlaceNameStateCity{PlaceName: pnsc.PlaceName, StateCity: pnsc.StateCity})
		}
	}()
	go func() {
		for _, dr := range drValues {
			drCh <- dr
		}
		close(drCh)
	}()
	fetchLocations(id, mockLogger, mockLocator, mockDataIO, &wg, drCh, pnscCh)
	close(pnscCh)
	_ = values

	for _, pnsc := range pnscValues {
		mockLogger.AssertCalled(t, "Debug", "goroutine", id, "getting location for", pnsc.PlaceName)
	}
	for _, pnsc := range pnscValues[:2] {
		mockLogger.AssertCalled(t, "Debug", "goroutine", id, "fetched location for", pnsc.PlaceName)
		mockLocator.AssertCalled(t, "GetLocation", pnsc.PlaceName)
	}
	err := fmt.Errorf("goroutine %d could not fetch location for '%s': %s", id, pnscValues[2].PlaceName, errGetLocation)
	mockLogger.AssertCalled(t, "Error", err)
	mockLogger.AssertCalled(t, "Debug", "goroutine", id, "done fetching locations")
}

func TestTranslateOK(t *testing.T) {
	mockLogger := newMockLogger()
	mockDataIO := newMockDataIO(false)
	mockLocator := newMockLocator()

	err := Translate(mockLogger, mockDataIO, mockLocator)

	assert.Nil(t, err)
	mockLogger.AssertCalled(t, "Info", "reading records")
	mockLogger.AssertCalled(t, "Info", "spawning goroutines to process data records")
	mockLogger.AssertCalled(t, "Info", "spawning goroutine to populate output csv")
	mockLogger.AssertCalled(t, "Debug", "closing dataRecordsChan")
	mockLogger.AssertCalled(t, "Debug", "closing placeNameStateCityChan")
	mockLogger.AssertCalled(t, "Debug", "waiting for all output records to be written")
	for _, pnsc := range pnscValues[:2] {
		mockDataIO.AssertCalled(t, "WritePlaceWithCity", pnsc.PlaceName, pnsc.StateCity)
	}
}

func TestTranslateReadError(t *testing.T) {
	mockLogger := newMockLogger()
	mockDataIO := newMockDataIO(true)
	mockLocator := newMockLocator()

	err := Translate(mockLogger, mockDataIO, mockLocator)

	errActual := fmt.Errorf("error translating records: %s", errReadRecords)
	assert.True(t, err.Error() == errActual.Error())
}

func newMockLogger() *logger.MockLogger {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Debug", mock.Anything).Return(nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Error", mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return mockLogger
}

func newMockDataIO(isReadRecordsError bool) *dataio.MockDataIO {
	mockDataIO := new(dataio.MockDataIO)
	if isReadRecordsError {
		mockDataIO.On("ReadRecords").Return(drValues, errReadRecords)
	} else {
		mockDataIO.On("ReadRecords").Return(drValues, nil)
	}
	mockDataIO.On("WritePlaceWithCity", mock.Anything, mock.Anything).Return(nil)
	return mockDataIO
}

func newMockLocator() *locator.MockLocator {
	mockLocator := new(locator.MockLocator)
	mockLocator.On("GetLocation", drValues[0].PlaceName).Return(&pnscValues[0].StateCity, nil)
	mockLocator.On("GetLocation", drValues[1].PlaceName).Return(&pnscValues[1].StateCity, nil)
	mockLocator.On("GetLocation", drValues[2].PlaceName).Return(&pnscValues[0].StateCity, errGetLocation)
	return mockLocator
}
