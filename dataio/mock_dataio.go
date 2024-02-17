package dataio

import (
	"github.com/SanferD/table-populator/domain"
	"github.com/stretchr/testify/mock"
)

type MockDataIO struct {
	mock.Mock
}

func (md *MockDataIO) ReadRecords() ([]domain.DataRecord, error) {
	args := md.Called(nil)
	return args.Get(0).([]domain.DataRecord), args.Error(1)
}

func (md *MockDataIO) WritePlaceWithCity(placeName string, stateCity domain.StateCity) error {
	args := md.Called(placeName, stateCity)
	return args.Error(0)
}
