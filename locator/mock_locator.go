package locator

import (
	"context"

	"github.com/SanferD/table-populator/domain"
	"github.com/stretchr/testify/mock"
	"googlemaps.github.io/maps"
)

type MockLocator struct {
	mock.Mock
}

func (ml *MockLocator) GetLocation(name string) (*domain.StateCity, error) {
	args := ml.Called(name)
	return args.Get(0).(*domain.StateCity), args.Error(1)
}

type MockMapsClient struct {
	mock.Mock
}

func (mmc *MockMapsClient) TextSearch(context context.Context, tsr *maps.TextSearchRequest) (maps.PlacesSearchResponse, error) {
	args := mmc.Called(context, tsr)
	return args.Get(0).(maps.PlacesSearchResponse), args.Error(1)
}
