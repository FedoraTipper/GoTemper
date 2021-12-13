package main

import (
	"fmt"
	"os"

	"github.com/FedoraTipper/gotemper/internal/driver"
)

func main() {
	devices, err := driver.FindTemperDevices()

	if err != nil {
		panic(err)
	}

	if len(devices) == 0 {
		fmt.Println("No Temper devices found")
		os.Exit(1)
	} else {
		fmt.Printf("%v\n", devices)
	}

	device := devices[len(devices)-1]
	deviceDriver := device.Details.Driver

	temp, err := deviceDriver.GetStats(device)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Temp: %f", temp)

}
