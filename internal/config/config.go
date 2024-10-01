package config

import "crypto/rsa"

// Режимы работы программы.
const (
	DBMode     = "DB"
	FileMode   = "file"
	MemoryMode = "memory"
)

// ServerOptions содержит настройки серверной части приложения.
var ServerOptions struct {
	Host             string
	FileStoragePath  string
	DBaddress        string
	Mode             string
	Key              string
	Restore          bool
	SignMode         bool
	StoreInterval    int64
	ReconnectCount   int
	ReconnectDelta   int
	PrivateKeyPath   string
	PrivateCryptoKey *rsa.PrivateKey
	ConfigJSONPath   string
	TrustedSubnet    string
	GrpcOn           bool
	GrpcHost         string
}

// ClientOptions содержит настройки серверной части приложения.
var ClientOptions struct {
	Host            string
	Key             string
	SignMode        bool
	RateLimit       int64
	PoolWorkers     int64
	ReportInterval  int64
	PollInterval    int64
	BatchSize       int64
	PublicKeyPath   string
	PrivateKeyPath  string
	PublicCryptoKey *rsa.PublicKey
	ConfigJSONPath  string
	GrpcOn          bool
	GrpcHost        string
}

// Переменные для хранения информации о версии сборки.
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// Тип для ведения информации о версии сборки.
type VersionInfo struct {
	BuildVersion string
	BuildCommit  string
	BuildDate    string
}

// SetMode устанавливает режим работы приложения, используется в тестировании.
func SetMode(mode string) {
	ServerOptions.Mode = mode
}

// GetVersionInfo возвращает информацию о версии сборки.
func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		BuildVersion: buildVersion,
		BuildCommit:  buildCommit,
		BuildDate:    buildDate,
	}
}

// SetSignMode устанавливает режим подписи отправляемыъ сообщений, используется в тестировании.
func SetSignMode(mode bool) {
	ClientOptions.SignMode = mode
}
