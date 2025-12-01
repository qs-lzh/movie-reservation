package main

import (
	"log"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/interfaces/web"
	"github.com/qs-lzh/movie-reservation/internal/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open gorm.DB: %v", err)
	}

	app := app.New(cfg, db)
	router := web.InitRouter(app)

	router.Run(cfg.Addr)

}
