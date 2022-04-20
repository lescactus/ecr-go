package appconfig

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env"
)

var (
	validLogLevels = []string{"error", "info", "debug"}
)

func LoadConfig(c *config) error {
	if err := env.Parse(&c.Application); err != nil {
		return err
	}
	if !isValidLogLevel(c.Application.LogLevel) {
		return errors.New("LogLevel must be 'error', 'info' or 'debug'")
	}
	return nil
}

func init() {
	Config = &config{}
	err := LoadConfig(Config)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse configuration:\n%v", err))
	}
}

func isValidLogLevel(l string) bool {
	for _, v := range validLogLevels {
		if v == l {
			return true
		}
	}
	return false
}
