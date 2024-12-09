package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	Url   string `yaml:"url" env:"VOGONIX_URL"`
	User  string `yaml:"user" env:"VOGONIX_USER"`
	Token string `yaml:"token" env:"VOGONIX_TOKEN"`
}

func ReadConfig(configFile string) (*ConfigDatabase, error) {

	var cfg ConfigDatabase
	cleanenv.ReadConfig(configFile, &cfg)
	cleanenv.ReadEnv(&cfg)
	return &cfg, nil
}
