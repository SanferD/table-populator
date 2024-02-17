package logger

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) Info(msgs ...any) {
	ml.Called(msgs...)
}

func (ml *MockLogger) Warn(msgs ...any) {
	ml.Called(msgs...)
}

func (ml *MockLogger) Debug(msgs ...any) {
	ml.Called(msgs...)
}

func (ml *MockLogger) Error(msgs ...any) {
	ml.Called(msgs...)
}

func (ml *MockLogger) Fatal(msgs ...any) {
	ml.Called(msgs...)
}
