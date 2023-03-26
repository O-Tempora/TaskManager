package apiserver

import "dip/internal/store"

type Config struct {
	Port     string        `yaml:"port"`
	LogLevel string        `yaml:"log_level"`
	Store    *store.Config `yaml:"store"`
}

func NewConfig() *Config {
	return &Config{
		Port:     ":5192",
		LogLevel: "Debug",
		Store:    store.NewConfig(),
	}
}
