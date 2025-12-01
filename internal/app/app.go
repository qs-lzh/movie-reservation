package app

import (
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type App struct {
	Config *config.Config

	DB *gorm.DB

	UserService        service.UserService
	MovieService       service.MovieService
	ShowtimeService    service.ShowtimeService
	ReservationService service.ReservationService
}

func New(config *config.Config, db *gorm.DB) *App {
	return &App{
		Config:             config,
		DB:                 db,
		UserService:        service.NewUserService(db),
		MovieService:       service.NewMovieService(db),
		ShowtimeService:    service.NewShowtimeService(db),
		ReservationService: service.NewReservationService(db),
	}
}

func (app *App) Close() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
