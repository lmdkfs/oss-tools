package config

type Config struct {
	Qiniu Qiniu
	Ufile Ufile
	Log   Log
	Worker int
}

type Log struct {
	LogName string
	LogPath string
}
type Qiniu struct {
	Ak            string
	Sk            string
	Bucket        string
	UpHost        string
	RsfHost       string
	RsHost        string
	IoVipHost     string
	ApiHost       string
	QiniuDomain   string
	PrivateBucket bool // true priviate false public
}
type Ufile struct {
	PublicKey  string
	PrivateKey string
	Bucket     string
	FileHost   string
	Prefix     string
}

var GlobalConfig *Config

func NewGlobalConfig() *Config {
	if GlobalConfig == nil {
		GlobalConfig = &Config{}
	}
	return GlobalConfig
}
