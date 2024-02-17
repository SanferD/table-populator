package locator

import (
	"github.com/SanferD/table-populator/domain"
	"github.com/stretchr/testify/mock"
)

type MockLocator struct {
	mock.Mock
}

func (ml *MockLocator) GetLocation(name string) (*domain.StateCity, error) {
	args := ml.Called(name)
	return args.Get(0).(*domain.StateCity), args.Error(1)
}
