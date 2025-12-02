package main

import (
	"log"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/interfaces/web"
	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/model"
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

	initDB(db)

	app := app.New(cfg, db)
	defer app.Close()
	router := web.InitRouter(app)

	router.Run(cfg.Addr)
}

func initDB(db *gorm.DB) {
	db.Migrator().AutoMigrate(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
	)
}
