package apiserver

type Config struct {
	Port     string    `yaml:"port"`
	LogLevel string    `yaml:"log_level"`
	DBconf   *DBconfig `yaml:"dbconfig"`
}

type DBconfig struct {
	Host     string `yaml:"host"`
	DBname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBort    string `yaml:"dbport"`
}

func NewConfig() *Config {
	return &Config{
		Port:     ":5192",
		LogLevel: "Debug",
	}
}
