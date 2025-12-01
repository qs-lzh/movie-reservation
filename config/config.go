package config

import (
	"os"

	"github.com/qs-lzh/movie-reservation/internal/util"
)

type Config struct {
	DatabaseDSN string
	Addr        string
}

func LoadConfig() (*Config, error) {
	if err := util.LoadEnv(); err != nil {
		return nil, err
	}
	databaseDSN := os.Getenv("DATABASE_DSN")
	addr := os.Getenv("ADDR")
	return &Config{
		DatabaseDSN: databaseDSN,
		Addr:        addr,
	}, nil
}
