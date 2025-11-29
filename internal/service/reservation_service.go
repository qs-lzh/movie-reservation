package service

import (
	"errors"

	"github.com/qs-lzh/movie-reservation/internal/repository"
	"gorm.io/gorm"
)

type ReservationService interface {
	Reserve(userName uint, showtimeID uint) error
}

type reservationService struct {
	db       *gorm.DB
	repo     repository.ReservationRepo
	userRepo repository.UserRepo
}

// var _ ReservationService = (reservationService)(nil)

func NewReservationService(db *gorm.DB) *reservationService {
	return &reservationService{
		db:       db,
		repo:     repository.NewReservationRepoGorm(db),
		userRepo: repository.NewUserRepoGorm(db),
	}
}

// func (s *reservationService) Reserve(userName string, showtimeID uint) error {
// 	// check if userID exists
// 	user, err := s.userRepo.GetByName(userName)
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return
// 	}
// 	if err != nil {
// 		return
// 	}
//
// 	// check if showtimeID exists
//
// 	// check if the ticket is stil available
//
// 	// reserve
// }
