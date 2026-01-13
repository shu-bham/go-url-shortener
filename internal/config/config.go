package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
	Logger struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logger"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(env string) (*Config, error) {
	if env == "" {
		env = "dev"
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	fileName := fmt.Sprintf("%s/../../configs/%s.yml", basepath, env)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &cfg, nil
}
