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
	AuthService        service.AuthService
}

func New(config *config.Config, db *gorm.DB) *App {
	userService := service.NewUserService(db)
	return &App{
		Config:             config,
		DB:                 db,
		UserService:        userService,
		MovieService:       service.NewMovieService(db),
		ShowtimeService:    service.NewShowtimeService(db),
		ReservationService: service.NewReservationService(db),
		AuthService:        service.NewJWTAuthService(userService),
	}
}

func (app *App) Close() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
