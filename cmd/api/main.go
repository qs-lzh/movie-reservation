package main

import (
	"log"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/interfaces/web"
	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/cache"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/security"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	security.InitJWT(cfg.JWTSecretKey)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open gorm.DB: %v", err)
	}
	initDB(db)

	cache := cache.NewRedisCache(cfg.CacheURL)

	app := app.New(cfg, db, cache)
	defer app.Close()
	router := web.InitRouter(app)

	if err := router.RunTLS(cfg.Addr, cfg.CertPath, cfg.KeyPath); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *gorm.DB) {
	db.Migrator().AutoMigrate(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
	)
}
