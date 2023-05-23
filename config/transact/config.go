package transact

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type GRPC struct {
	Addr string `yaml:"Addr"`
}

type DB struct {
	ConnString string `yaml:"ConnectionString"`
}

type App struct {
	Name    string `yaml:"Name"`
	Version string `yaml:"Version"`
}

type Auth struct {
	Secret string `yaml:"Secret"`
}

type Config struct {
	DB   `yaml:"DB"`
	App  `yaml:"App"`
	Auth `yaml:"Auth"`
	GRPC `yaml:"GRPC"`
}

var cfg Config

func NewConfig(filename string) (*Config, error) {
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can`t read config,exiting with error %w", err)
	}
	return &cfg, nil
}
