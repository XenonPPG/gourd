package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address         string        `yaml:"address"`
		ShutdownTimeout time.Duration `yaml:"shutdown-timeout"`
	} `yaml:"server"`
	DocumentService struct {
		WorkerCount int `yaml:"worker-count"`
	} `yaml:"document-service"`
}

func Load(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
