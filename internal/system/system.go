package system

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
)

type SystemDriver interface {
	Open(filePath string) error
	Write(payload []byte) error
	Read() ([]byte, error)
	Close()
}

func GetFileDriver(devPath string) (SystemDriver, error) {
	if compiledHIDRegex.MatchString(devPath) {
		return &HIDRawDriver{}, nil
	}

	return nil, errors.New(fmt.Sprintf("No system driver found for %s", devPath))
}

func FindDevPaths(devicePath string, devPathsMap map[string]string) ([]string, error) {
	entries, err := os.ReadDir(devicePath)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if compiledTTYRegex.MatchString(entry.Name()) || compiledHIDRegex.MatchString(entry.Name()) {
			devPathsMap[entry.Name()] = ""
		} else if entry.IsDir() {
			_, err := FindDevPaths(path.Join(devicePath, entry.Name()), devPathsMap)

			if err != nil {
				return nil, err
			}
		}
	}

	keys := make([]string, 0, len(devPathsMap))

	for key := range devPathsMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys, nil
}
