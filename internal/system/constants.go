package system

import "regexp"

const (
	devPath = "/dev"

	hidRegex = "hidraw[0-9]"
	ttyRegex = "tty.*[0-9]"
)

var (
	compiledHIDRegex = regexp.MustCompile(hidRegex)
	compiledTTYRegex = regexp.MustCompile(ttyRegex)
)
