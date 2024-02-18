package ioutil

import (
	"io"
	"os"

	"github.com/stretchr/testify/mock"
)

type MockFileOps struct {
	mock.Mock
}

func (miou *MockFileOps) Open(name string) (*os.File, error) {
	args := miou.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

func (miou *MockFileOps) Create(name string) (*os.File, error) {
	args := miou.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

func (miou *MockFileOps) Exit(id int) {
	_ = miou.Called(id)
}

type MockCSVReader struct {
	mock.Mock
}

func (mr *MockCSVReader) ReadAll() ([][]string, error) {
	args := mr.Called(nil)
	return args.Get(0).([][]string), args.Error(1)
}

type MockCSVWriter struct {
	mock.Mock
}

func (mw *MockCSVWriter) Write(row []string) error {
	args := mw.Called(row)
	return args.Error(0)
}

func (mw *MockCSVWriter) Flush() {
	_ = mw.Called()
}

type MockLogCreator struct {
	mock.Mock
}

func (mlc *MockLogCreator) New(w io.Writer, prefix string, flags int) Logger {
	args := mlc.Called(w, prefix, flags)
	return args.Get(0).(Logger)
}

type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) SetPrefix(prefix string) {
	_ = ml.Called(prefix)
}

func (ml *MockLogger) Println(msgs ...any) {
	_ = ml.Called(msgs...)
}

type MockCloser struct {
	mock.Mock
}

func (mc *MockCloser) Close() error {
	args := mc.Called()
	return args.Error(0)
}
