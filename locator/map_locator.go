package locator

import (
	"context"
	"fmt"
	"strings"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
	"googlemaps.github.io/maps"
)

type MapsClient interface {
	TextSearch(context.Context, *maps.TextSearchRequest) (maps.PlacesSearchResponse, error)
}

type GoogleMapsLocator struct {
	mapsClient MapsClient
}

func InitializeMapLocator(config config.Config) (*GoogleMapsLocator, error) {
	mapsClient, err := maps.NewClient(maps.WithAPIKey(config.MapsAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create new maps client: %s", err)
	}
	return &GoogleMapsLocator{mapsClient: mapsClient}, nil
}

func (mapsLocator *GoogleMapsLocator) GetLocation(placeName string) (*domain.StateCity, error) {
	mapsClient := mapsLocator.mapsClient

	// get maps places
	r := &maps.TextSearchRequest{
		Query: placeName,
	}
	res, err := mapsClient.TextSearch(context.Background(), r)
	if err != nil {
		return nil, fmt.Errorf("failed to get text response for '%s': %s", placeName, err)
	}

	// extract city and state code from formatted address
	if len(res.Results) == 0 {
		return nil, fmt.Errorf("failed to get any results for '%s'", placeName)
	}
	formattedAddress := res.Results[0].FormattedAddress
	addressParts := strings.Split(formattedAddress, ",")
	if len(addressParts) < 3 {
		err = fmt.Errorf("formatted address '%s' doesn't have 3 commas", formattedAddress)
		return nil, err
	}

	city := strings.TrimSpace(addressParts[len(addressParts)-3])
	stateZip := strings.TrimSpace(addressParts[len(addressParts)-2])
	state := strings.Fields(stateZip)[0]

	return &domain.StateCity{State: state, City: city}, nil
}
