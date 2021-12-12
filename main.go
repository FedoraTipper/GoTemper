package main

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

type USBID struct {
	VenderID  uint16
	ProductID uint16
	Offset    uint8
}

type TemperDevice struct {
	Path        string
	UsbID       USBID
	IsHIDRaw    bool
	DevicePaths []string
}

const (
	sysPath           = "/sys/bus/usb/devices"
	devPath           = "/dev"
	idVendorFileName  = "idVendor"
	idProductFileName = "idProduct"

	hidRegex = "hidraw[0-9]"
	ttyRegex = "tty.*[0-9]"
)

const (
	standardTemper = "temper"
	goldTemper     = "tempergold"
)

var (
	temperDevices = map[string]USBID{
		standardTemper: {
			VenderID:  0x0c45,
			ProductID: 0x7401,
			Offset:    2,
		},
		goldTemper: {
			VenderID:  0x413d,
			ProductID: 0x2107,
			Offset:    2,
		},
	}

	compiledHIDRegex = regexp.MustCompile(hidRegex)
	compiledTTYRegex = regexp.MustCompile(ttyRegex)
)

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

func findDevPaths(devicePath string, devPathsMap map[string]string) ([]string, error) {
	entries, err := os.ReadDir(devicePath)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if compiledTTYRegex.MatchString(entry.Name()) || compiledHIDRegex.MatchString(entry.Name()) {
			devPathsMap[entry.Name()] = ""
		} else if entry.IsDir() {
			_, err := findDevPaths(path.Join(devicePath, entry.Name()), devPathsMap)

			if err != nil {
				return nil, err
			}
		}
	}

	keys := make([]string, 0, len(devPathsMap))

	for key := range devPathsMap {
		keys = append(keys, key)
	}

	return keys, nil
}

func isTemperDevice(usbPath string) (*USBID, error) {
	vendorId, err := readIdFile(usbPath, idVendorFileName)

	if err != nil {
		return nil, err
	}

	productId, err := readIdFile(usbPath, idProductFileName)

	if err != nil {
		return nil, err
	}

	for device := range temperDevices {
		temper := temperDevices[device]

		if vendorId == temper.VenderID && productId == temper.ProductID {
			return &temper, nil
		}
	}

	return nil, nil
}

func findDevices() ([]TemperDevice, error) {
	devices := []TemperDevice{}

	walkDirFunc := func(path string, info fs.FileInfo, err error) error {
		temper, err := isTemperDevice(path)

		if err != nil {
			return err
		}

		if temper != nil {
			fmt.Println(path)
			fmt.Println("device found")

			devPaths, err := findDevPaths(path, map[string]string{})

			if err != nil {
				return err
			}

			devices = append(devices, TemperDevice{
				Path:        path,
				UsbID:       *temper,
				DevicePaths: devPaths,
			})
		}

		return nil
	}

	err := filepath.Walk(sysPath, walkDirFunc)

	if err != nil {
		return nil, err
	}

	return devices, nil
}

func readTemperature(device TemperDevice) (float64, error) {
	devicePath := device.DevicePaths[len(device.DevicePaths)-1]

	disk := path.Join(devPath, devicePath)
	fd, err := syscall.Open(disk, syscall.O_RDWR|syscall.O_NONBLOCK, 0777)

	defer func(fd int) {
		err := syscall.Close(fd)

		if err != nil {
			log.Println("Unable to close fd socket")
			log.Println(err)
		}

	}(fd)

	if err != nil {
		return 0, err
	}

	bytesWritten, err := syscall.Write(fd, []byte{0x01, 0x80, 0x33, 0x01, 0, 0, 0, 0})

	if err != nil {
		return 0, err
	}

	if bytesWritten != 8 {
		panic("issue")
	}

	buffer := make([]byte, 8, 8)
	for {
		bitsRead, err := syscall.Read(fd, buffer)

		if err != nil && err != syscall.EAGAIN {
			return 0, err
		}

		if bitsRead > 0 {
			break
		}
	}

	rawTempReading := binary.BigEndian.Uint16(buffer[device.UsbID.Offset:])
	temperature := float64(rawTempReading / 100)
	return temperature, err
}

func main() {
	devices, err := findDevices()

	if err != nil {
		panic(err)
	}

	if len(devices) == 0 {
		fmt.Println("No Temper devices found")
	} else {
		fmt.Printf("%v\n", devices)
	}

	device := devices[len(devices)-1]

	temp, err := readTemperature(device)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Temp: %f", temp)

}
