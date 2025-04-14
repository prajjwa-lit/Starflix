package config

// Config holds all the server configuration
type Config struct {
	Port          int
	VideoDir      string
	CoverImageDir string
	MaxUploadSize int // in MB
}

// NewConfig creates a default configuration
func NewConfig() *Config {
	return &Config{
		Port:          8080,
		VideoDir:      "./videos",
		CoverImageDir: "./covers",
		MaxUploadSize: 4096, // 1GB default
	}
}

// MaxUploadSizeBytes returns the max upload size in bytes
func (c *Config) MaxUploadSizeBytes() int64 {
	return int64(c.MaxUploadSize) * 1024 * 1024
}
