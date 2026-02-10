package config

import (
	"io"
)

var defaultDb = "binip.db"

type Config struct {
	DbFile      string
	DebugFile   string
	DebugWriter io.Writer
	Debug       bool
}

func NewConfig() *Config {
	return &Config{
		DbFile: defaultDb,
		Debug:  false,
	}
}
