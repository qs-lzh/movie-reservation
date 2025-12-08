package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type SeatService interface {
	CreateSeat(seat *model.Seat) error
	GetSeatByID(id uint) (*model.Seat, error)
	GetSeatsByHallID(hallID uint) ([]model.Seat, error)
	DeleteSeatByID(id uint) error
}

type seatService struct {
	db   *gorm.DB
	repo repository.SeatRepo
}

func NewseatService(db *gorm.DB) *seatService {
	return &seatService{
		db:   db,
		repo: repository.NewSeatRepoGorm(db),
	}
}

var _ SeatService = (*seatService)(nil)

func (s *seatService) CreateSeat(seat *model.Seat) error {
	if err := s.repo.Create(seat); err != nil {
		return err
	}
	return nil
}

func (s *seatService) GetSeatByID(id uint) (*model.Seat, error) {
	seat, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return seat, nil
}

func (s *seatService) DeleteSeatByID(id uint) error {
	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}
	return nil
}

func (s *seatService) GetSeatsByHallID(hallID uint) ([]model.Seat, error) {
	seats, err := s.repo.GetByHallID(hallID)
	if err != nil {
		return nil, err
	}
	return seats, nil
}
