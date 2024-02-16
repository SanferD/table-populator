package main

import (
	"fmt"

	"github.com/SanferD/table-populator/application"
	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/dataio"
	"github.com/SanferD/table-populator/locator"
	"github.com/SanferD/table-populator/logger"
)

func main() {
	// initialize config
	config, err := config.New()
	if err != nil {
		panic(fmt.Errorf("error initializing configuration: %s", err))
	}

	// initialize logger
	logger, err := logger.New(config)
	if err != nil {
		panic(fmt.Errorf("error initializing multi logger: %s", err))
	}

	// initialize dataio
	dataio, err := dataio.New(config)
	if err != nil {
		panic(fmt.Errorf("error creating dataio: %s", err))
	}

	// initialize locator
	locator, err := locator.New(config)
	if err != nil {
		panic(fmt.Errorf("error creating locator: %s", err))
	}

	// run application
	if err := application.Translate(logger, dataio, locator); err != nil {
		panic(fmt.Errorf("error creating translator: %s", err))
	}
}
