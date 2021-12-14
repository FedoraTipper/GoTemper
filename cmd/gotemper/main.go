package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FedoraTipper/gotemper/internal/output"
	"github.com/FedoraTipper/gotemper/internal/temper"
	"github.com/go-co-op/gocron"
)

func PostStats(devices []temper.TemperDevice, outputDriver output.OutputDriver) {
	device := devices[len(devices)-1]
	deviceDriver := device.Details.Driver

	temp, err := deviceDriver.GetStats(device)

	if err != nil {
		panic(err)
	}

	outputDriver.PostStats("Temp", fmt.Sprintf("%f", temp.InternalTemperature))
}

func main() {
	devices, err := temper.FindTemperDevices()

	if err != nil {
		panic(err)
	}

	if len(devices) == 0 {
		fmt.Println("No Temper devices found")
		os.Exit(1)
	} else {
		fmt.Printf("%v\n", devices)
	}

	outputDriver, err := output.GetOutputDriver("stdout")

	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	_, err = scheduler.Every("1s").Do(PostStats, devices, outputDriver)

	if err != nil {
		return
	}

	scheduler.StartBlocking()
}
