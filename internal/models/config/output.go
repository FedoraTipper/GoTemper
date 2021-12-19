package config

import (
	"errors"

	"github.com/FedoraTipper/gotemper/internal/constants"
)

type OutputDriverConfig struct {
	Output       string
	LoggingLevel string
	LoggingFile  string
	InfluxDB     InfluxDBConfig
}

func (o *OutputDriverConfig) Validate() []error {
	if len(o.Output) == 0 {
		return []error{errors.New("Output driver is not set in the configuration file")}
	}

	if constants.OutputType(o.Output) == constants.OutputTypeInfluxDb {
		return o.InfluxDB.Validate()
	}

	return nil
}
