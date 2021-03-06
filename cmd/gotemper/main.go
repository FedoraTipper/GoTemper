package main

import (
	"fmt"
	"log"
	"time"

	"github.com/FedoraTipper/gotemper/internal/config"
	"github.com/FedoraTipper/gotemper/internal/logger"
	configModels "github.com/FedoraTipper/gotemper/internal/models/config"
	"github.com/FedoraTipper/gotemper/internal/output"
	"github.com/FedoraTipper/gotemper/internal/temper"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

var (
	configName  = "config.yml"
	configPaths = []string{
		".",
		"./configs/",
		"$HOME/gotemper/",
		"$HOME/.config/gotemper/",
	}
)

var defaultConfigValues = map[string]interface{}{
	"LoggingLevel": "info",
	"LoggingFile":  "",
	"Output":       "stdout",
}

func PostStats(devices []temper.TemperDevice, outputDriver output.OutputDriver) {
	device := devices[len(devices)-1]
	deviceDriver := device.Details.Driver

	temp, err := deviceDriver.GetStats(device)

	if err != nil {
		panic(err)
	}

	outputDriver.PostStats("temperature", "internal", temp.InternalTemperature)
}

func main() {
	viper := config.GenerateConfigReader(defaultConfigValues, configName, configPaths)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error when reading in config.yml")
		log.Fatalf("%v", err)
	}

	var outputConfig configModels.OutputDriverConfig

	if err := viper.UnmarshalExact(&outputConfig); err != nil {
		log.Println("Error in parsing config.yml")
		log.Fatalf("%v", err)
	}

	configErrs := outputConfig.Validate()

	if len(configErrs) > 0 {
		log.Fatalf("%v", configErrs)
	}

	err := logger.InitGlobalLogger(outputConfig.LoggingFile, outputConfig.LoggingLevel)

	if err != nil {
		log.Fatalf("Unable to configure logger. Error: %v", err)
	}

	outputDriver, err := output.GetOutputDriver(outputConfig.Output)

	if err != nil {
		zap.S().Fatalf("%v", err)
	}

	err = outputDriver.Initialise(outputConfig)

	if err != nil {
		zap.S().Fatalf("%v", err)
	}

	devices, err := temper.FindTemperDevices()

	if err != nil {
		zap.S().Fatalf("%v", err)
		log.Fatalf("%v", err)
	}

	if len(devices) == 0 {
		msg := "No Temper devices found"
		fmt.Println(msg)
		zap.S().Fatal(msg)
	} else {
		fmt.Printf("%v\n", devices)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	_, err = scheduler.Every("1s").Do(PostStats, devices, outputDriver)

	if err != nil {
		zap.S().Fatalf("%v", err)
	}

	scheduler.StartBlocking()
}
