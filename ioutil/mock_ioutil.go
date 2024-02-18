package ioutil

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockIOUtil struct {
	mock.Mock
}

func (miou *MockIOUtil) Open(name string) (*os.File, error) {
	args := miou.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

func (miou *MockIOUtil) Create(name string) (*os.File, error) {
	args := miou.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

type MockCSVReader struct {
	mock.Mock
}

type MockCSVWriter struct {
	mock.Mock
}

func (mr *MockCSVReader) ReadAll() ([][]string, error) {
	args := mr.Called(nil)
	return args.Get(0).([][]string), args.Error(1)
}

func (mw *MockCSVWriter) Write(row []string) error {
	args := mw.Called(row)
	return args.Error(0)
}

func (mw *MockCSVWriter) Flush() {
	_ = mw.Called(nil)
}
