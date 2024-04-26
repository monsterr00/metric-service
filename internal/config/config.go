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
}

var ClientOptions struct {
	Host           string
	ReportInterval int64
	PollInterval   int64
	BatchSize      int64
}

const (
	DBMode     = "DB"
	FileMode   = "file"
	MemoryMode = "memory"
)
