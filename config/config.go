package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	AlertBot AlertBot `yaml:"alert_bot"`
}

type AlertBot struct {
	Token string `yaml:"telegram_token"`
	Host  string `yaml:"telegram_host"`
}

type Config struct {
	Server `yaml:"server"`
}

func New(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	d := yaml.NewDecoder(bytes.NewReader(content))
	if err = d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
