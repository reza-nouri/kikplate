package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v2"
)

type CLIConfig struct {
	Server ServerConfig `yaml:"server"`
	Auth   AuthConfig   `yaml:"auth"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type AuthConfig struct {
	Token string `yaml:"token"`
}

var DefaultConfigPath = filepath.Join(os.Getenv("HOME"), ".kikplate", "config.yaml")

func Load(path string) (*CLIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %w", path, err)
	}
	var cfg CLIConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}
	if cfg.Server.Address == "" {
		return nil, fmt.Errorf("server.address is required in config file")
	}
	return &cfg, nil
}

func Save(path string, cfg *CLIConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
