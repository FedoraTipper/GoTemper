package output

import (
	"errors"
	"fmt"
)

type OutputDriver interface {
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
