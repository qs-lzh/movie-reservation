package config

import (
	"os"

	"github.com/qs-lzh/movie-reservation/internal/util"
)

type Config struct {
	DatabaseDSN  string
	Addr         string
	JWTSecretKey string
	CertPath     string
	KeyPath      string
}

func LoadConfig() (*Config, error) {
	if err := util.LoadEnv(); err != nil {
		return nil, err
	}
	databaseDSN := os.Getenv("DATABASE_DSN")
	addr := os.Getenv("ADDR")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	crtPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")
	return &Config{
		DatabaseDSN:  databaseDSN,
		Addr:         addr,
		JWTSecretKey: jwtSecretKey,
		CertPath:     crtPath,
		KeyPath:      keyPath,
	}, nil
}
