package output

import "fmt"

type StdOutDriver struct {
}

func (s *StdOutDriver) PostStats(label, payload string) {
	fmt.Printf("%s: %s", label, payload)
}
