package main

import (
	"fmt"
	loggerMod "github.com/SanferD/table-populator/logger"
	configMod "github.com/SanferD/table-populator/config"
	dataioMod "github.com/SanferD/table-populator/dataio"
	locatorMod "github.com/SanferD/table-populator/locator"
	"github.com/SanferD/table-populator/application"
)

func main() {
	// initialize config
	config, err := configMod.InitializeConfig()
	if err != nil {
		panic(fmt.Errorf("error initializing configuration: %s", err))
	}

	// initialize logger
	logger, err := loggerMod.CreateLogger(config)
	if err != nil {
		panic(fmt.Errorf("error initializing multi logger: %s", err))
	}

	// initialize dataio
	dataio, err := dataioMod.CreateDataIo(config)
	if err != nil {
		panic(fmt.Errorf("error creating dataio: %s", err))
	}

	// initialize locator
	locator, err := locatorMod.CreateLocator(config)
	if err != nil {
		panic(fmt.Errorf("error creating locator: %s", err))
	}

	// run application
	if err := application.Translate(logger, dataio, locator); err != nil {
		panic(fmt.Errorf("error creating translator: %s", err))
	}
}
