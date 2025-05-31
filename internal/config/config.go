package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Add your configuration fields here
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}
