package output

import (
	"context"
	"errors"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDBDriver struct {
	client *influxdb2.Client
	api    *api.WriteAPI
}

func (i *InfluxDBDriver) Initialise(outputConfig OutputDriverConfig) error {
	config := outputConfig.InfluxDB
	client := influxdb2.NewClient(config.ServerAddress, config.Token)
	writeAPI := client.WriteAPI(config.Org, config.Bucket)

	i.api = &writeAPI
	i.client = &client

	resp, err := client.Ping(context.Background())

	if err != nil {
		return err
	}

	if !resp {
		return errors.New(fmt.Sprintf("Unable to ping InfluxDB instance %s", config.ServerAddress))
	}

	return nil
}

func (i *InfluxDBDriver) PostStats(label, payload string) {
	panic("implement me")
}
