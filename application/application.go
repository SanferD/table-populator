package application

import (
	"time"
	"sync"
	"fmt"
	"github.com/SanferD/table-populator/domain"
)

const numWorkerThreads = 1
const getLocationRateLimit = time.Second / 5 // 5 times per second

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

	// buffered queue to hold the records for the threads to work on
	dataRecordsChan := make(chan domain.DataRecord, numWorkerThreads)
	placeNameStateCityChan := make(chan PlaceNameStateCity, numWorkerThreads)
	var wg sync.WaitGroup

	// spawn multiple threads to process the data records in the queue
	logger.Info("spawning threads to process data records")
	for i := 0; i<numWorkerThreads; i++ {
		wg.Add(1)
		go fetchLocations(logger, locationGetter, dataIo, &wg, dataRecordsChan, placeNameStateCityChan)
	}

	// spawn thread to write placeName, state, city to output.csv
	logger.Info("spawning thread to populate output csv")
	wg.Add(1)
	go populateOutput(logger, dataIo, &wg, placeNameStateCityChan)

	// add the records to the queue for the worker threads to process
	for _, record := range records {
		dataRecordsChan <- record
	}
	
	// wait for the worker threads to complete
	close(dataRecordsChan)
	close(placeNameStateCityChan)
	wg.Wait()
	return nil
}

func populateOutput(logger domain.Logger, dataIo domain.DataIo, wg *sync.WaitGroup,
					placeNameStateCityChan chan PlaceNameStateCity) {
	defer wg.Done()
	for placeNameStateCity := range placeNameStateCityChan {
		logger.Debug("writing place with city", placeNameStateCity.PlaceName, placeNameStateCity.StateCity)
		dataIo.WritePlaceWithCity(placeNameStateCity.PlaceName, placeNameStateCity.StateCity)
	}
}

func fetchLocations(logger domain.Logger, locationGetter domain.LocationGetter, dataIo domain.DataIo,
					wg *sync.WaitGroup, dataRecordChan chan domain.DataRecord,
					dataRecordWithStateCityChan chan PlaceNameStateCity) {
	defer wg.Done()
	throttle := time.Tick(getLocationRateLimit)

	for dataRecord := range dataRecordChan {
		<- throttle
		// get location
		logger.Debug("getting location for", dataRecord.PlaceName)
		stateCity, err := locationGetter.GetLocation(dataRecord.PlaceName)
		if err != nil {
			err = fmt.Errorf("could not fetch location for '%s': %s", dataRecord.PlaceName, err)
			logger.Error(err)
			continue
		}
		dataRecordWithStateCityChan <- PlaceNameStateCity{
			PlaceName: dataRecord.PlaceName,
			StateCity: *stateCity,}
	}
}