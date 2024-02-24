package util

type UtilConfig struct {
	CacheDir string
	FontDir  string
}

var (
	Config = &UtilConfig{}
)

func SetUtilConfig(config *UtilConfig) {
	Config = config
}
