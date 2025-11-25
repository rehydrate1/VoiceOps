package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SSH        SSHConfig        `yaml:"ssh"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Commands   []Command        `yaml:"commands"`
}

type SSHConfig struct {
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	KeyPath string `yaml:"key_path"`
}

type MonitoringConfig struct {
	URLs []string `yaml:"urls"`
}

type Command struct {
	Phrase   string `yaml:"phrase"`
	Script   string `yaml:"script"`
	Response string `yaml:"response"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
