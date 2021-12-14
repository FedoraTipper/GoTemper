package output

type OutputType string

const (
	OutputTypeStdOut   OutputType = "stdout"
	OutputTypeInfluxDb OutputType = "influxdb"
)
