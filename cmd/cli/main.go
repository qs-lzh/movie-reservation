package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

func main() {
	dsn := "host=localhost user=liuzihao password=cvbn123.. dbname=movie_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
