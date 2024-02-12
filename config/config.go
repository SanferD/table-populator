package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	CsvDataFilePath   string  `mapstructure:"CSV_DATA_FILE_PATH"`	
	MapsApiKey 		  string  `mapstructure:"MAPS_API_KEY"`
	OutputCsvFilePath string  `mapstructure:"OUTPUT_FILE_PATH"`
	DataIoKind		  string  `mapstructure:"DATAIO_KIND"`
	LocatorKind 	  string  `mapstructure:"LOCATOR_KIND"`
	LoggerKind 		  string  `mapstructure:"LOGGER_KIND"`
	LogToStdout 	  bool 	  `mapstructure:"LOG_TO_STDOUT"`
	LogFilePath 	  *string `mapstructure:"LOG_FILE_PATH"`
}

func InitializeConfig() (Config, error) {
	var config Config

	// configure
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName(".env")

	// set defaults
	v.SetDefault("DATAIO_KIND", "csv")
	v.SetDefault("LOCATOR_KIND", "google-maps")
	v.SetDefault("LOGGER_KIND", "multi")
	v.SetDefault("LOG_TO_STDOUT", true)
	v.SetDefault("LOG_FILE_PATH", nil)

	// read config
	if err := v.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading in configuration: %s", err)
	}

	// parse config
	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("error unmarshalling configuration: %s", err)
	}
	return config, nil
}
