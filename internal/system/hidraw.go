package system

import (
	"errors"
	"fmt"
	"log"
	"path"
	"syscall"

	"github.com/FedoraTipper/gotemper/internal/constants"
)

const (
	fileRWPerm = 0644
)

type HIDRawDriver struct {
	fd int
}

func (h *HIDRawDriver) Open(filePath string) error {
	filePath = path.Join(constants.DevPath, filePath)
	fd, err := syscall.Open(filePath, syscall.O_RDWR|syscall.O_NONBLOCK, fileRWPerm)

	if err != nil {
		return err
	}

	h.fd = fd
	return nil
}

func (h *HIDRawDriver) Write(payload []byte) error {
	bytesWritten, err := syscall.Write(h.fd, payload)

	if err != nil {
		return err
	}

	if bytesWritten != len(payload) {
		return errors.New(fmt.Sprintf("Length of bytes written (%d) does not match payload length (%d)", bytesWritten, len(payload)))
	}

	return nil
}

func (h *HIDRawDriver) Read() ([]byte, error) {
	buffer := make([]byte, 8, 8)

	for {
		bytesRead, err := syscall.Read(h.fd, buffer)

		// Skip resource temporarily unavailable error
		if err != nil && err != syscall.EAGAIN {
			return nil, err
		}

		if bytesRead > 0 {
			break
		}
	}

	return buffer, nil
}

func (h *HIDRawDriver) Close() {
	err := syscall.Close(h.fd)

	if err != nil {
		log.Printf("Unable to close file descriptor %d", h.fd)
		log.Println(err)
	}
}
