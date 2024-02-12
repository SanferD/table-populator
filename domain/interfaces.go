package domain

type DataIo interface {
	ReadRecords() ([]DataRecord, error)
	WritePlaceWithCity(placeName string, stateCity StateCity) error
}

type LocationGetter interface {
	GetLocation(name string) (*StateCity, error)
}

type Logger interface {
	Info(msgs ...any)
	Warn(msgs ...any)
	Debug(msgs ...any)
	Error(msgs ...any)
	Fatal(msgs ...any)
}