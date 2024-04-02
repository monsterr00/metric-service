package config

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

var ServerOptions struct {
	Host string
}

var ClientOptions struct {
	Host           string
	ReportInterval int64
	PollInterval   int64
}
