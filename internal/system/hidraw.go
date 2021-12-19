package system

import (
	"errors"
	"fmt"
	"path"
	"syscall"

	"github.com/FedoraTipper/gotemper/internal/constants"
	"go.uber.org/zap"
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

	zap.S().Infow("Successfully opened file descriptor", "fd", h.fd, "File Path", filePath)

	h.fd = fd
	return nil
}

func (h *HIDRawDriver) Write(payload []byte) error {
	zap.S().Debugw("Beginning syscall write", "payload length (bytes)", len(payload))
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

	zap.S().Debugf("Beginning syscall read of %d", h.fd)
	for {
		bytesRead, err := syscall.Read(h.fd, buffer)

		// Skip resource temporarily unavailable error
		if err != nil && err != syscall.EAGAIN {
			return nil, err
		}

		if bytesRead > 0 {
			zap.S().Debugf("Bytes read from syscall (len: %d)", bytesRead)
			break
		}
	}

	return buffer, nil
}

func (h *HIDRawDriver) Close() {
	err := syscall.Close(h.fd)

	if err != nil {
		zap.S().Errorw("Unable to close file descriptor", "fd", h.fd, "Error", err)
	}
}
