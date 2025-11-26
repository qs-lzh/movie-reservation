package app

import (
	"database/sql"

	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type App struct {
	DB *sql.DB
	MovieRepo *
}
