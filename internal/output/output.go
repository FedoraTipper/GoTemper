package output

import (
	"errors"
	"fmt"

	"github.com/FedoraTipper/gotemper/internal/constants"
	"github.com/FedoraTipper/gotemper/internal/models/config"
)

type OutputDriver interface {
	Initialise(config config.OutputDriverConfig) error
	PostStats(label, sublabel string, payload interface{})
}

func GetOutputDriver(typ string) (OutputDriver, error) {
	outputType := constants.OutputType(typ)

	switch outputType {
	case constants.OutputTypeStdOut:
		return &StdOutDriver{}, nil
	case constants.OutputTypeInfluxDb:
		return &InfluxDBDriver{}, nil
	}

	return nil, errors.New(fmt.Sprintf("Unable to find output driver for %s", typ))
}
