package config

import (
	"github.com/FedoraTipper/gotemper/internal/constants"
)

type OutputDriverConfig struct {
	Output       string
	LoggingLevel string
	LoggingFile  string
	InfluxDB     InfluxDBConfig
}

func (o *OutputDriverConfig) Validate() []error {
	if constants.OutputType(o.Output) == constants.OutputTypeInfluxDb {
		return o.InfluxDB.Validate()
	}

	return nil
}
