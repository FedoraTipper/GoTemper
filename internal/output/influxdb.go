package output

import (
	"context"
	"errors"
	"fmt"

	configModels "github.com/FedoraTipper/gotemper/internal/models/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"go.uber.org/zap"
)

type InfluxDBDriver struct {
	client *influxdb2.Client
	api    *api.WriteAPI
}

func (i *InfluxDBDriver) Initialise(outputConfig configModels.OutputDriverConfig) error {
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

	zap.S().Debugw("InfluxDB output initialised")

	return nil
}

func (i *InfluxDBDriver) PostStats(label, sublabel string, payload interface{}) {
	api := *i.api

	go func() {
		zap.S().Debugw("Posting stats", "label", label, "sublabel", sublabel, "payload", payload)
		point := influxdb2.NewPointWithMeasurement(label).AddField(sublabel, payload)
		api.WritePoint(point)
		api.Flush()
	}()

	select {
	case err := <-api.Errors():
		zap.S().Errorw("Error when posting stats to influxdb", "Error", err)
		break
	}
}
