package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type ShowtimeSeatService interface {
	CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error
	GetShowtimeSeatByID(id uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatByShowtimeID(showtimdID uint) ([]model.ShowtimeSeat, error)
	GetShowtimeSeatByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error)
	UpdateShowtimeSeatStatus(id uint, targetStatus model.ShowtimeSeatStatus) error
	DeleteShowtimeSeatByID(id uint) error
}

type showtimeSeatService struct {
	db   *gorm.DB
	repo repository.ShowtimeSeatRepo
}

func NewshowtimeSeatService(db *gorm.DB) *showtimeSeatService {
	return &showtimeSeatService{
		db:   db,
		repo: repository.NewshowtimeSeatRepoGorm(db),
	}
}

var _ ShowtimeSeatService = (*showtimeSeatService)(nil)

func (s *showtimeSeatService) CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error {
	if err := s.repo.Create(showtimeSeat); err != nil {
		return err
	}
	return nil
}

func (s *showtimeSeatService) GetShowtimeSeatByID(id uint) (*model.ShowtimeSeat, error) {
	showtimeSeat, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return showtimeSeat, nil
}

func (s *showtimeSeatService) GetShowtimeSeatByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error) {
	return s.repo.GetByShowtimeID(showtimeID)
}

func (s *showtimeSeatService) GetShowtimeSeatByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error) {
	return s.repo.GetByStatus(status)
}

var ErrShowtimeSeatNotExist = errors.New("The ShowtimeSeat does not exist")
var ErrShowtimeSeatStatusNotChange = errors.New("The target status of showtimeSeat is the same as the origin")

func (s *showtimeSeatService) UpdateShowtimeSeatStatus(id uint, targetStatus model.ShowtimeSeatStatus) error {
	existingShowtimeSeat, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrShowtimeSeatNotExist
		}
		return err
	}

	if targetStatus == existingShowtimeSeat.Status {
		return ErrShowtimeSeatStatusNotChange
	}

	existingShowtimeSeat.Status = targetStatus
	return s.repo.Update(id, existingShowtimeSeat)
}

func (s *showtimeSeatService) DeleteShowtimeSeatByID(id uint) error {
	return s.repo.DeleteByID(id)
}
