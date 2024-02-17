package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/SanferD/table-populator/domain"
)

const numWorkerGoroutines = 3
const getLocationRateLimit = time.Second / 2 // 2 times per second

type PlaceNameStateCity struct {
	PlaceName string
	StateCity domain.StateCity
}

func newDataRecordChan() chan domain.DataRecord {
	// buffered channel to hold the records for the location getter goroutines to translate placeName to location
	return make(chan domain.DataRecord, numWorkerGoroutines)
}

func newPlaceNameStateCityChan() chan PlaceNameStateCity {
	// buffered channel to hold the records for the placeName goroutines to write (placeName, stateCity) to output csv
	return make(chan PlaceNameStateCity, numWorkerGoroutines)
}

func newDonePopulateOutputChannel() chan bool {
	// channel to notify when all the output records have been written to the output
	return make(chan bool)
}

func Translate(logger domain.Logger, dataIO domain.DataIO, locator domain.Locator) error {
	// read records from file
	logger.Info("reading records")
	records, err := dataIO.ReadRecords()
	if err != nil {
		return fmt.Errorf("error translating records: %s", err)
	}

	drCh := newDataRecordChan()
	pnscCh := newPlaceNameStateCityChan()
	dpoCh := newDonePopulateOutputChannel()
	var wg sync.WaitGroup // used to block the top-level function until all the location getter goroutines have terminated

	// spawn multiple goroutines to process the data records in the queue
	logger.Info("spawning goroutines to process data records")
	for i := 0; i < numWorkerGoroutines; i++ {
		wg.Add(1)
		go fetchLocations(i, logger, locator, dataIO, &wg, drCh, pnscCh)
	}

	// spawn goroutine to write placeName, state, city to output.csv
	logger.Info("spawning goroutine to populate output csv")
	go populateOutput(logger, dataIO, pnscCh, dpoCh)

	// add the records to the queue for the worker goroutine to process
	for _, record := range records {
		drCh <- record
	}

	logger.Debug("closing dataRecordsChan")
	close(drCh)
	wg.Wait()
	logger.Debug("closing placeNameStateCityChan")
	close(pnscCh)
	logger.Debug("waiting for all output records to be written")
	<-dpoCh
	return nil
}

func populateOutput(logger domain.Logger, dataIO domain.DataIO, pnscCh chan PlaceNameStateCity, dpoCh chan bool) {
	for pnsc := range pnscCh {
		logger.Debug("writing place with city", pnsc.PlaceName, pnsc.StateCity)
		dataIO.WritePlaceWithCity(pnsc.PlaceName, pnsc.StateCity)
	}
	dpoCh <- true
}

func fetchLocations(id int, logger domain.Logger, locator domain.Locator, dataIO domain.DataIO, wg *sync.WaitGroup, drCh chan domain.DataRecord, pnscCh chan PlaceNameStateCity) {
	defer wg.Done()
	rateLimiter := time.NewTicker(getLocationRateLimit)

	for dataRecord := range drCh {
		<-rateLimiter.C

		// get location
		logger.Debug("goroutine", id, "getting location for", dataRecord.PlaceName)
		stateCity, err := locator.GetLocation(dataRecord.PlaceName)
		if err != nil {
			err = fmt.Errorf("goroutine %d could not fetch location for '%s': %s", id, dataRecord.PlaceName, err)
			logger.Error(err)
			continue
		}
		logger.Debug("goroutine", id, "fetched location for", dataRecord.PlaceName)
		pnscCh <- PlaceNameStateCity{
			PlaceName: dataRecord.PlaceName,
			StateCity: *stateCity}
	}
	logger.Debug("goroutine", id, "done fetching locations")
}
