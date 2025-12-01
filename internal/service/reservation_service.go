package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ReservationService interface {
	Reserve(userID uint, showtimeID uint) error
	CancelReservation(reservationID uint) error
	GetRemainingTickets(showtimeID uint) (int, error)
}

type reservationService struct {
	db           *gorm.DB
	repo         repository.ReservationRepo
	showtimeRepo repository.ShowtimeRepo
	userRepo     repository.UserRepo
}

var _ ReservationService = (*reservationService)(nil)

func NewReservationService(db *gorm.DB) *reservationService {
	return &reservationService{
		db:           db,
		repo:         repository.NewReservationRepoGorm(db),
		showtimeRepo: repository.NewShowtimeRepoGorm(db),
		userRepo:     repository.NewUserRepoGorm(db),
	}
}

// NOTE  logic wrong!
func (s *reservationService) Reserve(userID uint, showtimeID uint) error {
	// check if showtime exists
	_, err := s.showtimeRepo.GetByID(showtimeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrShowtimeNotExist
		}
		return err
	}

	// check if the ticket is stil available
	count, err := s.GetRemainingTickets(showtimeID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if count <= 0 {
		return ErrNoTicketsAvailable
	}

	// check if the user already have the same reservation
	reservations, err := s.repo.GetByUserID(userID)
	for _, reservation := range reservations {
		if reservation.ShowtimeID == showtimeID {
			return ErrAlreadyReserved
		}
	}

	// reserve
	return s.repo.Create(&model.Reservation{
		ShowtimeID: showtimeID,
		UserID:     userID,
	})
}

func (s *reservationService) CancelReservation(reservationID uint) error {
	if _, err := s.repo.GetByID(reservationID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.DeleteByID(reservationID)
}

func (s *reservationService) GetRemainingTickets(showtimeID uint) (int, error) {
	reservations, err := s.repo.GetByShowtimeID(showtimeID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return model.SeatCount, nil
		}
		return 0, err
	}
	return model.SeatCount - len(reservations), nil
}
