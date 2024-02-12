package application

import (
	"sync"
	"fmt"
	"github.com/SanferD/table-populator/domain"
)

const NUM_WORKER_THREADS = 2

func Translate(logger domain.Logger, dataIo domain.DataIo,
			   locationGetter domain.LocationGetter) error {
	// read records from file
	logger.Info("reading records")
	records, err := dataIo.ReadRecords()
	if err != nil {
		return fmt.Errorf("error translating records: %s", err)
	}

	// buffered queue to hold the records for the threads to work on
	dataRecordsQueue := make(chan domain.DataRecord, NUM_WORKER_THREADS)
	var wg sync.WaitGroup

	// spawn multiple threads to process the data records in the queue
	logger.Info("translating records")
	for i := 0; i<NUM_WORKER_THREADS; i++ {
		wg.Add(1)
		go processDataRecords(logger, locationGetter, dataIo, &wg, dataRecordsQueue)
	}

	// add the records to the queue for the worker threads to process
	for _, record := range records {
		dataRecordsQueue <- record
	}
	
	// wait for the worker threads to complete
	close(dataRecordsQueue)
	wg.Wait()
	return nil
}

func processDataRecords(logger domain.Logger, locationGetter domain.LocationGetter, dataIo domain.DataIo,
						wg *sync.WaitGroup, dataRecordQueue chan domain.DataRecord) {
	defer wg.Done()

	for dataRecord := range dataRecordQueue {
		// get location
		logger.Debug("getting location for", dataRecord.PlaceName)
		stateCity, err := locationGetter.GetLocation(dataRecord.PlaceName)
		if err != nil {
			err = fmt.Errorf("could not fetch location for '%s': %s", dataRecord.PlaceName, err)
			logger.Error(err)
			continue
		}

		// write location to csv
		logger.Debug("writing place with city", dataRecord.PlaceName, stateCity.State, stateCity.City)
		dataIo.WritePlaceWithCity(dataRecord.PlaceName, *stateCity)
	}
}