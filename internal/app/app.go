package app

import (
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/internal/cache"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type App struct {
	Config *config.Config

	DB    *gorm.DB
	Cache *cache.RedisCache

	UserService        service.UserService
	MovieService       service.MovieService
	ShowtimeService    service.ShowtimeService
	ReservationService service.ReservationService
	AuthService        service.AuthService
	CaptchaService     service.CaptchaService
}

func New(config *config.Config, db *gorm.DB, cache *cache.RedisCache) *App {
	userService := service.NewUserService(db)
	return &App{
		Config:             config,
		DB:                 db,
		Cache:              cache,
		UserService:        userService,
		MovieService:       service.NewMovieService(db),
		ShowtimeService:    service.NewShowtimeService(db),
		ReservationService: service.NewReservationService(db),
		AuthService:        service.NewJWTAuthService(userService),
		CaptchaService:     service.NewCaptchaService(cache),
	}
}

func (app *App) Close() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
