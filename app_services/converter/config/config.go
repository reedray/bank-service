package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Grpc struct {
	Addr string `yaml:"Addr"`
}

type Redis struct {
	Addr      string        `yaml:"Addr"`
	ExpiresAt time.Duration `yaml:"expiresAt"`
}

type Config struct {
	LogLevel string `yaml:"LogLevel"`
	App      `yaml:"App"`
	Grpc     `yaml:"GRPC"`
	Redis    `yaml:"Redis"`
}

var cfg Config

func NewConfig(filename string) (*Config, error) {
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can`t read config,exiting with error %w", err)
	}
	return &cfg, nil
}
