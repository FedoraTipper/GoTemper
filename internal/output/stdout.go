package output

import (
	"fmt"

	configModels "github.com/FedoraTipper/gotemper/internal/models/config"
)

type StdOutDriver struct {
}

func (s *StdOutDriver) Initialise(config configModels.OutputDriverConfig) error {
	return nil
}

func (s *StdOutDriver) PostStats(label, sublabel string, payload interface{}) {
	fmt.Printf("%s-%s: %s", label, sublabel, payload)
}
