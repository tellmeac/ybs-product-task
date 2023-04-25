package core

import (
	"github.com/spf13/viper"
	"yandex-team.ru/bstask/internal/pkg/web"
	"yandex-team.ru/bstask/internal/storage"
)

type Config struct {
	RPS     int              `yaml:"rps"`
	Server  web.ServerConfig `yaml:"server"`
	Storage storage.Config   `yaml:"storage"`
}

func ParseConfig(loader *viper.Viper) (*Config, error) {
	cfg := &Config{}

	if err := loader.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := loader.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
