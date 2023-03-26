package store

type Config struct {
	Host     string `yaml:"host"`
	DBname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewConfig() *Config {
	return &Config{}
}
