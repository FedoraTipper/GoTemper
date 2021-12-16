package output

import (
	"errors"
	"fmt"
)

type OutputDriverConfig struct {
	InfluxDB InfluxDBConfig
}

type InfluxDBConfig struct {
	ServerAddress string
	Token         string
	Bucket        string
	Org           string
}

type OutputDriver interface {
	Initialise(config OutputDriverConfig) error
	PostStats(label, payload string)
}

func GetOutputDriver(typ string) (OutputDriver, error) {
	outputType := OutputType(typ)

	switch outputType {
	case OutputTypeStdOut:
		return &StdOutDriver{}, nil
	case OutputTypeInfluxDb:
		return &InfluxDBDriver{}, nil
	}

	return nil, errors.New(fmt.Sprintf("Unable to find output driver for %s", typ))
}
