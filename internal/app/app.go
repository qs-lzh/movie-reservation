package app

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/internal/cache"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type App struct {
	Config *config.Config

	DB     *gorm.DB
	Cache  *cache.RedisCache
	Logger *zap.Logger

	UserService         service.UserService
	MovieService        service.MovieService
	ShowtimeService     service.ShowtimeService
	ReservationService  service.ReservationService
	HallService         service.HallService
	ShowtimeSeatService service.ShowtimeSeatService
	AuthService         service.AuthService
	CaptchaService      service.CaptchaService
}

func New(config *config.Config, db *gorm.DB, cache *cache.RedisCache, logger *zap.Logger) *App {
	userService := service.NewUserService(db)
	return &App{
		Config:              config,
		DB:                  db,
		Cache:               cache,
		Logger:              logger,
		UserService:         userService,
		MovieService:        service.NewMovieService(db),
		ShowtimeService:     service.NewShowtimeService(db),
		ReservationService:  service.NewReservationService(db),
		HallService:         service.NewHallService(db),
		ShowtimeSeatService: service.NewshowtimeSeatService(db),
		AuthService:         service.NewJWTAuthService(userService),
		CaptchaService:      service.NewCaptchaService(cache),
	}
}

func (app *App) Close() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
