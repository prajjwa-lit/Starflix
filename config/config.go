package config

type Config struct {
	Port          int
	VideoDir      string
	CoverImageDir string
	MaxUploadSize int
}

func NewConfig() *Config {
	return &Config{
		Port:          8080,
		VideoDir:      "./videos",
		CoverImageDir: "./covers",
		MaxUploadSize: 4096,
	}
}

func (c *Config) MaxUploadSizeBytes() int64 {
	return int64(c.MaxUploadSize) * 1024 * 1024
}
