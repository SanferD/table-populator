package ioutil

import (
	"io"
	"log"
	"os"
)

type FileOps interface {
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	Exit(code int)
}

type StdFileOps struct{}

func (o *StdFileOps) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (o *StdFileOps) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (o *StdFileOps) Exit(code int) {
	os.Exit(code)
}

type CSVReader interface {
	ReadAll() ([][]string, error)
}

type CSVWriter interface {
	Write([]string) error
	Flush()
}

type Logger interface {
	SetPrefix(string)
	Println(...any)
}

type LogCreator interface {
	New(io.Writer, string, int) Logger
}

type StdLogCreator struct{}

func (StdLogCreator) New(out io.Writer, prefix string, flag int) Logger {
	return log.New(out, prefix, flag)
}
