package constants

import "regexp"

const (
	DevPath = "/dev"

	hidRegex = "hidraw[0-9]"
	ttyRegex = "tty.*[0-9]"
)

var (
	CompiledHIDRegex = regexp.MustCompile(hidRegex)
	CompiledTTYRegex = regexp.MustCompile(ttyRegex)
)
