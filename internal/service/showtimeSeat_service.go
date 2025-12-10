package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

/*
* ShowtimeSeatService do not examine if the operand is allowed or is safe
 */

type ShowtimeSeatService interface {
	CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error
	InitShowtimeSeatsForShowtime(showtime *model.Showtime) error
	GetShowtimeSeatByID(id uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatByShowtimeIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error)
	GetShowtimeSeatsByShowtimeID(showtimdID uint) ([]model.ShowtimeSeat, error)
	GetShowtimeSeatsByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error)
	UpdateShowtimeSeatStatus(id uint, targetStatus model.ShowtimeSeatStatus) error
	DeleteShowtimeSeatByID(id uint) error
}

type showtimeSeatService struct {
	db          *gorm.DB
	repo        repository.ShowtimeSeatRepo
	seatService SeatService
}

func NewShowtimeSeatService(db *gorm.DB, showtimeSeatRepo repository.ShowtimeSeatRepo, seatService SeatService) *showtimeSeatService {
	return &showtimeSeatService{
		db:          db,
		repo:        showtimeSeatRepo,
		seatService: seatService,
	}
}

var _ ShowtimeSeatService = (*showtimeSeatService)(nil)

func (s *showtimeSeatService) CreateShowtimeSeat(showtimeSeat *model.ShowtimeSeat) error {
	if err := s.repo.Create(showtimeSeat); err != nil {
		return err
	}
	return nil
}

func (s *showtimeSeatService) InitShowtimeSeatsForShowtime(showtime *model.Showtime) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		hallID := showtime.HallID
		seats, err := s.seatService.GetSeatsByHallID(hallID)
		if err != nil {
			return err
		}
		var showtimeSeats []model.ShowtimeSeat
		for _, seat := range seats {
			showtimeSeats = append(showtimeSeats, model.ShowtimeSeat{
				ShowtimeID: showtime.ID,
				SeatID:     seat.ID,
				Status:     model.StatusAvailable,
			})
		}
		return s.repo.CreateBatch(showtimeSeats)
	})
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

func (s *showtimeSeatService) GetShowtimeSeatByShowtimeIDSeatID(showtimeID, seatID uint) (*model.ShowtimeSeat, error) {
	showtimeSeat, err := s.repo.GetByShowIDSeatID(showtimeID, seatID)
	if err != nil {
		return nil, err
	}
	return showtimeSeat, nil
}

func (s *showtimeSeatService) GetShowtimeSeatsByShowtimeID(showtimeID uint) ([]model.ShowtimeSeat, error) {
	return s.repo.GetByShowtimeID(showtimeID)
}

func (s *showtimeSeatService) GetShowtimeSeatsByStatus(status model.ShowtimeSeatStatus) ([]model.ShowtimeSeat, error) {
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
