package locator

import (
	"fmt"

	"github.com/SanferD/table-populator/config"
	"github.com/SanferD/table-populator/domain"
)

type LocatorKind int

const (
	Error      LocatorKind = -1
	GoogleMaps LocatorKind = iota
)

func extractLocatorKind(kind string) (LocatorKind, error) {
	switch kind {
	case "google-maps":
		return GoogleMaps, nil
	default:
		return Error, fmt.Errorf("unrecognized locator kind string '%s'", kind)
	}
}

func New(config config.Config) (domain.Locator, error) {
	locatorKind, err := extractLocatorKind(config.LocatorKind)
	switch locatorKind {
	case GoogleMaps:
		return InitializeMapLocator(config)
	case Error:
		fallthrough
	default:
		return nil, fmt.Errorf("error extracting locator kind: %s", err)
	}
}
