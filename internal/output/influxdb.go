package output

type InfluxDBDriver struct {
}

func (i *InfluxDBDriver) PostStats(label, payload string) {
	panic("implement me")
}
