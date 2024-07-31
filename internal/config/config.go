package config

var ServerOptions struct {
	Host            string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	DBaddress       string
	Mode            string
	ReconnectCount  int
	ReconnectDelta  int
	Key             string
	SignMode        bool
}

var ClientOptions struct {
	Host           string
	ReportInterval int64
	PollInterval   int64
	Key            string
	BatchSize      int64
	SignMode       bool
	RateLimit      int64
	PoolWorkers    int64
}

const (
	DBMode     = "DB"
	FileMode   = "file"
	MemoryMode = "memory"
)

func SetMode(mode string) {
	ServerOptions.Mode = mode
}
