package temper

import (
	"encoding/binary"

	"github.com/FedoraTipper/gotemper/internal/system"
)

var (
	offset             = 2
	tempSysCallPayload = []byte{0x01, 0x80, 0x33, 0x01, 0, 0, 0, 0}
)

type StandardDriver struct {
}

func (t *StandardDriver) readTemperature(device TemperDevice) (float64, error) {
	devicePath := device.DevicePaths[len(device.DevicePaths)-1]

	fileDriver, err := system.GetFileDriver(devicePath)

	if err != nil {
		return 0, err
	}

	err = fileDriver.Open(devicePath)

	if err != nil {
		return 0, err
	}

	defer fileDriver.Close()

	err = fileDriver.Write(tempSysCallPayload)

	if err != nil {
		return 0, err
	}

	buffer, err := fileDriver.Read()

	if err != nil {
		return 0, err
	}

	rawTempReading := binary.BigEndian.Uint16(buffer[device.Details.Offset:])
	temperature := float64(rawTempReading / 100)
	return temperature, err
}

func (t *StandardDriver) GetStats(device TemperDevice) (TemperStats, error) {
	var temperStats TemperStats

	temp, err := t.readTemperature(device)

	if err != nil {
		return temperStats, err
	}

	temperStats.InternalTemperature = temp

	return temperStats, nil
}
