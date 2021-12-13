package driver

const (
	sysPath = "/sys/bus/usb/devices"

	idVendorFileName  = "idVendor"
	idProductFileName = "idProduct"
)

const (
	standardTemper = "temper"
	goldTemper     = "tempergold"
)

type DeviceDetails struct {
	Path      string
	VenderID  uint16
	ProductID uint16
	Offset    int
	Driver    TemperDriver
}

type TemperDevice struct {
	Details     DeviceDetails
	DevicePaths []string
}

var (
	temperDevicesMap = map[string]DeviceDetails{
		standardTemper: {
			VenderID:  0x0c45,
			ProductID: 0x7401,
			Offset:    2,
			Driver:    &StandardDriver{},
		},
		goldTemper: {
			VenderID:  0x413d,
			ProductID: 0x2107,
			Offset:    2,
			Driver:    &StandardDriver{},
		},
	}
)
