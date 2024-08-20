package config

const (
	DBMode     = "DB"
	FileMode   = "file"
	MemoryMode = "memory"
)

// ServerOptions содержит настройки серверной части приложения.
var ServerOptions struct {
	Host            string
	FileStoragePath string
	DBaddress       string
	Mode            string
	Key             string
	Restore         bool
	SignMode        bool
	StoreInterval   int64
	ReconnectCount  int
	ReconnectDelta  int
}

// ClientOptions содержит настройки серверной части приложения.
var ClientOptions struct {
	Host           string
	Key            string
	SignMode       bool
	RateLimit      int64
	PoolWorkers    int64
	ReportInterval int64
	PollInterval   int64
	BatchSize      int64
}

// SetMode устанавливает режим работы приложения, используется в тестировании.
func SetMode(mode string) {
	ServerOptions.Mode = mode
}
