package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	Url   string `yaml:"url" env:"URL"`
	User  string `yaml:"user" env:"USER"`
	Token string `yaml:"token" env:"TOKEN"`
}

func ReadConfig(configFile string) (*ConfigDatabase, error) {

	var cfg ConfigDatabase
	cleanenv.ReadConfig(configFile, &cfg)
	cleanenv.ReadEnv(&cfg)
	return &cfg, nil
}
