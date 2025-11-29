package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func main() {
	err := util.LoadEnv()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	db.AutoMigrate(
		model.Movie{},
		model.Showtime{},
	)
}
