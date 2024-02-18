package ioutil

import (
	"os"
)

type FileOpener interface {
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
}

type OSFileOpener struct{}

func (o OSFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (o OSFileOpener) Create(name string) (*os.File, error) {
	return os.Create(name)
}

type CSVReader interface {
	ReadAll() ([][]string, error)
}

type CSVWriter interface {
	Write([]string) error
	Flush()
}
