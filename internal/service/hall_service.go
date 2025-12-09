package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type HallService interface {
	CreateHall(hall *model.Hall) error
	UpdateHall(hall *model.Hall) error
	DeleteHallByID(id uint) error
	GetHallByID(id uint) (*model.Hall, error)
	GetHallByName(name string) (*model.Hall, error)
	GetAllHalls() ([]model.Hall, error)
}

type hallService struct {
	db          *gorm.DB
	repo        repository.HallRepo
	seatService SeatService
}

var _ HallService = (*hallService)(nil)

func NewHallService(db *gorm.DB, hallRepo repository.HallRepo, seatService SeatService) *hallService {
	return &hallService{
		db:          db,
		repo:        hallRepo,
		seatService: seatService,
	}
}

func (s *hallService) CreateHall(hall *model.Hall) error {
	if err := s.repo.Create(hall); err != nil {
		return err
	}
	return s.seatService.InitSeatsForHall(hall)
}

func (s *hallService) UpdateHall(hall *model.Hall) error {
	// NOTE: Ensure no related Seat and ShowtimeSeat

	// Verify that the hall with this ID exists
	existinghall, err := s.repo.GetByID(uint(hall.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// check if the new title is already used by another
	// because the title needs to be unique
	if existinghall.Name != hall.Name {
		anotherhall, err := s.repo.GetByName(hall.Name)
		if err == nil && anotherhall != nil && anotherhall.ID != hall.ID {
			return ErrAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return s.repo.Update(hall)
}

func (s *hallService) DeleteHallByID(id uint) error {
	// NOTE: Ensure no related Showtime and Seat and ShowtimeSeat

	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}
	return nil
}

func (s *hallService) GetHallByID(id uint) (*model.Hall, error) {
	hall, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return hall, nil
}

func (s *hallService) GetHallByName(name string) (*model.Hall, error) {
	hall, err := s.repo.GetByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return hall, nil
}

func (s *hallService) GetAllHalls() ([]model.Hall, error) {
	halls, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return halls, nil
}
