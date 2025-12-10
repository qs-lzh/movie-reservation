package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type ReservationService interface {
	// WARNING: need transaction
	Reserve(userID, showtimeID, seatID uint) error
	// WARNING: need transaction
	CancelReservation(reservationID uint) error
	// WARNING: need transaction
	GetRemainingTickets(showtimeID uint) (int, error)
	GetReservationsByUserID(userID uint) ([]model.Reservation, error)
	GetReservationByID(reservationID uint) (*model.Reservation, error)
}

type reservationService struct {
	db                  *gorm.DB
	repo                repository.ReservationRepo
	showtimeRepo        repository.ShowtimeRepo
	hallRepo            repository.HallRepo
	showtimeSeatService ShowtimeSeatService
}

var _ ReservationService = (*reservationService)(nil)

func NewReservationService(db *gorm.DB, reservationRepo repository.ReservationRepo,
	showtimeRepo repository.ShowtimeRepo, hallRepo repository.HallRepo, showtimeSeatService ShowtimeSeatService) *reservationService {
	return &reservationService{
		db:                  db,
		repo:                reservationRepo,
		showtimeRepo:        showtimeRepo,
		hallRepo:            hallRepo,
		showtimeSeatService: showtimeSeatService,
	}
}

func (s *reservationService) Reserve(userID, showtimeID, seatID uint) error {
	s.db.Transaction(func(tx *gorm.DB) error {
		// check if showtime exists
		showtime, err := s.showtimeRepo.GetByID(showtimeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrShowtimeNotExist
			}
			return err
		}

		// check if the ticket is stil available
		reservations, err := s.repo.GetByShowtimeID(showtimeID)
		if err != nil {
			// gorm.Find returns nil error for no records, so this is a real error
			return err
		}
		hall, err := s.hallRepo.GetByID(showtime.HallID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		count := hall.SeatCount - len(reservations)
		if count <= 0 {
			return ErrNoTicketsAvailable
		}

		return hall.SeatCount - len(reservations), nil
		if count <= 0 {
			return ErrNoTicketsAvailable
		}

		// check if the user already have the same reservation
		reservations, err := s.repo.GetByUserID(userID)
		if err != nil {
			return err
		}
		for _, reservation := range reservations {
			if reservation.ShowtimeID == showtimeID {
				return ErrAlreadyReserved
			}
		}

		// reserve
		if err = s.repo.Create(&model.Reservation{
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			UserID:     userID,
		}); err != nil {
			return err
		}

		// change showtimeSeat status
		showtimeSeat, err := s.showtimeSeatService.GetShowtimeSeatByShowtimeIDSeatID(showtimeID, seatID)
		if err != nil {
			return err
		}
		if err := s.showtimeSeatService.UpdateShowtimeSeatStatus(showtimeSeat.ID, model.StatusSold); err != nil {
			return err
		}
		return nil

	})
}

func (s *reservationService) ensureTicketAvailable() {

}

func (s *reservationService) CancelReservation(reservationID uint) error {
	reservation, err := s.repo.GetByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := s.repo.DeleteByID(reservationID); err != nil {
		return err
	}

	// change showtimeSeat status
	showtimeSeat, err := s.showtimeSeatService.GetShowtimeSeatByShowtimeIDSeatID(reservation.ShowtimeID, reservation.SeatID)
	if err != nil {
		return err
	}
	if err := s.showtimeSeatService.UpdateShowtimeSeatStatus(showtimeSeat.ID, model.StatusSold); err != nil {
		return err
	}
	return nil
}

func (s *reservationService) GetRemainingTickets(showtimeID uint) (int, error) {
	reservations, err := s.repo.GetByShowtimeID(showtimeID)
	if err != nil {
		// gorm.Find returns nil error for no records, so this is a real error
		return 0, err
	}
	// get total seats
	showtime, err := s.showtimeRepo.GetByID(showtimeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, nil
	}
	s.hallRepo.showtime.HallID

	return model.DefaultSeatCount - len(reservations), nil
}

func (s *reservationService) GetReservationsByUserID(userID uint) ([]model.Reservation, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		// gorm.Find returns nil error for no records, so this is a real error
		return nil, err
	}
	return reservations, nil
}

func (s *reservationService) GetReservationByID(reservationID uint) (*model.Reservation, error) {
	reservation, err := s.repo.GetByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return reservation, nil
}
