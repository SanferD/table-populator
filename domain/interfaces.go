package domain

type DataIO interface {
	ReadRecords() ([]DataRecord, error)
	WritePlaceWithCity(placeName string, stateCity StateCity) error
}

type Locator interface {
	GetLocation(name string) (*StateCity, error)
}

type Logger interface {
	Info(msgs ...any)
	Warn(msgs ...any)
	Debug(msgs ...any)
	Error(msgs ...any)
	Fatal(msgs ...any)
}
