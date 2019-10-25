package config

type Config struct {
	EndPoint        string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
	FileName        string

	BucketName string
	Location   string
	ObjectName string
}

var config *Config

// config
func NewConfig() *Config {
	if config == nil {
		config = &Config{}
	}

	return config
}
