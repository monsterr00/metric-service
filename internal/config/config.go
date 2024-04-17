package config

var ServerOptions struct {
	Host            string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
}

var ClientOptions struct {
	Host           string
	ReportInterval int64
	PollInterval   int64
}
