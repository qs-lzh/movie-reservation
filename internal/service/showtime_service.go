package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ShowtimeService interface {
	CreateShowtime(movieID uint, startTime time.Time, hallID uint) error
	UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error
	DeleteShowtime(showtimeID uint) error
	GetShowtimeByID(showtimeID uint) (*model.Showtime, error)
	ListShowtimesByMovieID(movieID uint) ([]model.Showtime, error)
}

type showtimeService struct {
	db   *gorm.DB
	repo repository.ShowtimeRepo
}

var _ ShowtimeService = (*showtimeService)(nil)

func NewShowtimeService(db *gorm.DB) *showtimeService {
	return &showtimeService{
		db:   db,
		repo: repository.NewShowtimeRepoGorm(db),
	}
}

func (s *showtimeService) CreateShowtime(movieID uint, startTime time.Time, hallID uint) error {
	showtime := &model.Showtime{
		MovieID: uint(movieID),
		StartAt: startTime,
		HallID:  uint(hallID),
	}
	if err := s.repo.Create(showtime); err != nil {
		return err
	}
	return nil
}

func (s *showtimeService) UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error {
	showtime, err := s.repo.GetByID(uint(showtimeID))
	if err != nil {
		return err
	}
	if err := s.repo.DeleteByID(uint(showtimeID)); err != nil {
		return err
	}
	showtime.ID = uint(showtimeID)
	showtime.StartAt = startTime
	showtime.HallID = uint(hallID)
	if err := s.repo.Create(showtime); err != nil {
		return err
	}
	return nil
}

func (s *showtimeService) DeleteShowtime(showtimeID uint) error {
	if err := s.repo.DeleteByID(uint(showtimeID)); err != nil {
		return err
	}
	return nil
}

func (s *showtimeService) GetShowtimeByID(showtimeID uint) (*model.Showtime, error) {
	showtime, err := s.repo.GetByID(uint(showtimeID))
	if err != nil {
		return &model.Showtime{}, err
	}
	return showtime, nil
}

func (s *showtimeService) ListShowtimesByMovieID(movieID uint) ([]model.Showtime, error) {
	showtimes, err := s.repo.GetByMovieID(uint(movieID))
	if err != nil {
		return nil, err
	}
	return showtimes, nil
}
