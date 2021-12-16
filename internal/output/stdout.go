package output

import "fmt"

type StdOutDriver struct {
}

func (s *StdOutDriver) Initialise(config OutputDriverConfig) error {
	return nil
}

func (s *StdOutDriver) PostStats(label, payload string) {
	fmt.Printf("%s: %s", label, payload)
}
