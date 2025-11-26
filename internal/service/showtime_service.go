package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ShowtimeService interface {
	CreateShowtime(movieID int, startTime time.Time, hallID int) error
	UpdateShowtime(showtimeID int, startTime time.Time, hallID int) error
	DeleteShowtime(showtimeID int) error
	GetShowtimeByID(showtimeID int) (*model.Showtime, error)
	ListShowtimesByMovie(movieID int) ([]model.Showtime, error)
	ListShowtimesByDate(date time.Time) ([]model.Showtime, error)
}

type showtimeService struct {
	db   *gorm.DB
	repo repository.ShowtimeRepo
}

func NewShowtimeService(db *gorm.DB) *showtimeService {
	return &showtimeService{
		db:   db,
		repo: repository.NewShowtimeRepoGorm(db),
	}
}

func (s *showtimeService) CreateShowtime(movieID int, startTime time.Time, hallID int) error {
	showtime := &model.Showtime{
		MovieID: uint(movieID),
		StartAt: startTime,
		HallID:  uint(hallID),
	}
	if err := s.repo.Create(*showtime); err != nil {
		return err
	}
	return nil
}

//
// func (s *showtimeService) UpdateShowtime(showtimeID int, startTime time.Time, hallID int) error {
//
// }
//
// func (s *showtimeService) DeleteShowtime(showtimeID int) error {
//
// }
//
// func (s *showtimeService) GetShowtimeByID(showtimeID int) (*model.Showtime, error) {
//
// }
//
// func (s *showtimeService) ListShowtimesByMovie(movieID int) ([]model.Showtime, error) {
//
// }
//
// func (s *showtimeService) ListShowtimesByDate(date time.Time) ([]model.Showtime, error) {
//
// }
