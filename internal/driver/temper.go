package driver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/FedoraTipper/gotemper/internal/system"
)

type TemperStats struct {
	InternalTemperature float64
}

type TemperDriver interface {
	GetStats(device TemperDevice) (TemperStats, error)
}

func readFileAndClean(filepath string) (string, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Printf("File %s doesn't exist, returning nil err to skip", filepath)
		return "", nil
	}

	byteData, err := os.ReadFile(filepath)

	if err != nil {
		return "", err
	}

	contents := string(byteData)

	// clean
	contents = strings.Trim(contents, "\n")

	return contents, nil
}

func readIdFile(usbPath, idFile string) (uint16, error) {
	idFile = path.Join(usbPath, idFile)
	contents, err := readFileAndClean(idFile)

	if err != nil {
		return 0, err
	}

	if len(contents) == 0 {
		return 0, nil
	}

	id, err := strconv.ParseUint(contents, 16, 64)

	return uint16(id), err
}

func GetUSBDetails(usbPath string) (*DeviceDetails, error) {
	vendorId, err := readIdFile(usbPath, idVendorFileName)

	if err != nil {
		return nil, err
	}

	productId, err := readIdFile(usbPath, idProductFileName)

	if err != nil {
		return nil, err
	}

	return &DeviceDetails{
		Path:      usbPath,
		VenderID:  vendorId,
		ProductID: productId,
	}, nil
}

func (details *DeviceDetails) addTemperDriver() error {
	for _, temperDevice := range temperDevicesMap {
		if details.VenderID == temperDevice.VenderID && details.ProductID == temperDevice.ProductID {
			driver := temperDevice.Driver
			details.Driver = driver
			details.Offset = temperDevice.Offset
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Device (VID: %d; PID: %d) was not found in preconfigured temper map", details.VenderID, details.ProductID))
}

func FindTemperDevices() ([]TemperDevice, error) {
	devices := []TemperDevice{}

	entries, err := os.ReadDir(sysPath)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Type() == os.ModeSymlink || entry.IsDir() {
			devPath := path.Join(sysPath, entry.Name())
			usbDetails, err := GetUSBDetails(devPath)

			if err != nil {
				return nil, err
			}

			err = usbDetails.addTemperDriver()

			if err != nil {
				log.Println(err) // should be warning
				log.Println("Skipping...")
				continue
			}

			devPaths, err := system.FindDevPaths(devPath, map[string]string{})

			if err != nil {
				return nil, err
			}

			devices = append(devices, TemperDevice{
				Details:     *usbDetails,
				DevicePaths: devPaths,
			})
		}
	}

	return devices, nil
}
