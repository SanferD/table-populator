package application

import (
	"time"
	"sync"
	"fmt"
	"github.com/SanferD/table-populator/domain"
)

const numWorkerGoroutines = 3
const getLocationRateLimit = time.Second / 2 // 2 times per second

type PlaceNameStateCity struct {
	PlaceName 	string
	StateCity   domain.StateCity
}

func Translate(logger domain.Logger, dataIo domain.DataIo,
			   locationGetter domain.LocationGetter) error {
	// read records from file
	logger.Info("reading records")
	records, err := dataIo.ReadRecords()
	if err != nil {
		return fmt.Errorf("error translating records: %s", err)
	}

	// buffered channel to hold the records for the location getter goroutines to translate placeName to location
	dataRecordsChan := make(chan domain.DataRecord, numWorkerGoroutines)
	// buffered channel to hold the records for the placeName goroutines to write (placeName, stateCity) to output csv
	placeNameStateCityChan := make(chan PlaceNameStateCity, numWorkerGoroutines)
	// channel to notify when all the output records have been written to the output
	donePopulatingOutputChan := make(chan bool)
	var wg sync.WaitGroup // used to block the top-level function until all the location getter goroutines have terminated

	// spawn multiple goroutines to process the data records in the queue
	logger.Info("spawning goroutines to process data records")
	for i := 0; i<numWorkerGoroutines; i++ {
		wg.Add(1)
		go fetchLocations(i, logger, locationGetter, dataIo, &wg, dataRecordsChan, placeNameStateCityChan)
	}

	// spawn goroutine to write placeName, state, city to output.csv
	logger.Info("spawning goroutine to populate output csv")
	go populateOutput(logger, dataIo, &wg, placeNameStateCityChan, donePopulatingOutputChan)

	// add the records to the queue for the worker goroutine to process
	for _, record := range records {
		dataRecordsChan <- record
	}
	
	logger.Debug("closing dataRecordsChan")
	close(dataRecordsChan)
	wg.Wait()
	logger.Debug("closing placeNameStateCityChan")
	close(placeNameStateCityChan)
	logger.Debug("waiting for all output records to be written")
	<- donePopulatingOutputChan
	return nil
}

func populateOutput(logger domain.Logger, dataIo domain.DataIo,
					placeNameStateCityChan chan PlaceNameStateCity,
					donePopulatingOutputChan chan bool) {
	for placeNameStateCity := range placeNameStateCityChan {
		logger.Debug("writing place with city", placeNameStateCity.PlaceName, placeNameStateCity.StateCity)
		dataIo.WritePlaceWithCity(placeNameStateCity.PlaceName, placeNameStateCity.StateCity)
	}
	donePopulatingOutputChan <- true
}

func fetchLocations(id int, logger domain.Logger, locationGetter domain.LocationGetter,
					dataIo domain.DataIo, wg *sync.WaitGroup,
					dataRecordChan chan domain.DataRecord,
					dataRecordWithStateCityChan chan PlaceNameStateCity) {
	defer wg.Done()
	throttle := time.Tick(getLocationRateLimit)

	for dataRecord := range dataRecordChan {
		<- throttle
		// get location
		logger.Debug("goroutine", id, "getting location for", dataRecord.PlaceName)
		stateCity, err := locationGetter.GetLocation(dataRecord.PlaceName)
		if err != nil {
			err = fmt.Errorf("could not fetch location for '%s': %s", dataRecord.PlaceName, err)
			logger.Error(err)
			continue
		}
		logger.Debug("goroutine", id, "fetched location for", dataRecord.PlaceName)
		dataRecordWithStateCityChan <- PlaceNameStateCity{
			PlaceName: dataRecord.PlaceName,
			StateCity: *stateCity,}
	}
	logger.Debug("goroutine", id, "done fetching locations")
}