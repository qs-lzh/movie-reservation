package app

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/internal/cache"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type App struct {
	Config *config.Config

	DB     *gorm.DB
	Cache  *cache.RedisCache
	Logger *zap.Logger

	UserRepo         *repository.UserRepo
	MovieRepo        *repository.MovieRepo
	ShowtimeRepo     *repository.ShowtimeRepo
	ReservationRepo  *repository.ReservationRepo
	HallRepo         *repository.HallRepo
	SeatRepo         *repository.SeatRepo
	ShowtimeSeatRepo *repository.ShowtimeSeatRepo

	UserService         service.UserService
	MovieService        service.MovieService
	ShowtimeService     service.ShowtimeService
	ReservationService  service.ReservationService
	HallService         service.HallService
	SeatService         service.SeatService
	ShowtimeSeatService service.ShowtimeSeatService
	AuthService         service.AuthService
	CaptchaService      service.CaptchaService
}

func New(config *config.Config, db *gorm.DB, cache *cache.RedisCache, logger *zap.Logger) *App {

	userRepo := repository.NewUserRepoGorm(db)
	movieRepo := repository.NewMovieRepoGorm(db)
	showtimeRepo := repository.NewShowtimeRepoGorm(db)
	reservationRepo := repository.NewReservationRepoGorm(db)
	hallRepo := repository.NewHallRepoGorm(db)
	seatRepo := repository.NewSeatRepoGorm(db)
	showtimeSeatRepo := repository.NewShowtimeSeatRepoGorm(db)

	userService := service.NewUserService(db, userRepo)
	seatService := service.NewseatService(db, seatRepo)
	showtimeSeatService := service.NewShowtimeSeatService(db, showtimeSeatRepo, seatService)
	showtimeService := service.NewShowtimeService(db, showtimeRepo, showtimeSeatService)
	reservationService := service.NewReservationService(db, reservationRepo, showtimeRepo, showtimeSeatService)
	captchaService := service.NewCaptchaService(cache)

	return &App{
		Config:              config,
		DB:                  db,
		Cache:               cache,
		Logger:              logger,
		UserService:         userService,
		MovieService:        service.NewMovieService(db, movieRepo, showtimeSeatService),
		ShowtimeService:     showtimeService,
		ReservationService:  reservationService,
		HallService:         service.NewHallService(db, hallRepo, seatService),
		SeatService:         seatService,
		ShowtimeSeatService: showtimeSeatService,
		AuthService:         service.NewJWTAuthService(userService),
		CaptchaService:      captchaService,
	}
}

func (app *App) Close() error {
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
