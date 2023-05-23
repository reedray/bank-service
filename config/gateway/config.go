package gateway

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Hosts struct {
	Gateway   string `yaml:"Gateway"`
	Converter string `yaml:"Converter"`
	Transact  string `yaml:"Transact"`
	Credit    string `yaml:"Credit"`
}

type App struct {
	Name    string `yaml:"Name"`
	Version string `yaml:"Version"`
}

type Config struct {
	Hosts `yaml:"Hosts"`
	App   `yaml:"App"`
}

var cfg Config

func NewConfig(filename string) (*Config, error) {
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can`t read config,exiting with error %w", err)
	}
	return &cfg, nil
}
