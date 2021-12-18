package config

import "errors"

type InfluxDBConfig struct {
	ServerAddress string
	Token         string
	Bucket        string
	Org           string
}

func (i *InfluxDBConfig) Validate() []error {
	var errs []error

	if i.ServerAddress == "" {
		errs = append(errs, errors.New("ServerAddress is empty in InfluxDB config"))
	}

	if i.Token == "" {
		errs = append(errs, errors.New("Token is empty in InfluxDB config"))
	}

	if i.Bucket == "" {
		errs = append(errs, errors.New("Bucket is empty in InfluxDB config"))
	}

	if i.Org == "" {
		errs = append(errs, errors.New("Org is empty in InfluxDB config"))
	}

	return errs
}
