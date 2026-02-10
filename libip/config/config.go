package config

import (
	"io"
)

type Config struct {
	DebugFile   string
	DebugWriter io.Writer
	Debug       bool
}

func NewConfig() *Config {
	return &Config{
		Debug: false,
	}
}
