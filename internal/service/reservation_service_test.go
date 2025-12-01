package service_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
	"github.com/qs-lzh/movie-reservation/internal/util"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestReservationService(t *testing.T) {
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	// drop and recreate tables
	db.Migrator().DropTable(
		&model.Movie{},
		&model.Showtime{},
		&model.User{},
		&model.Reservation{},
	)
	db.AutoMigrate(
		&model.Movie{},
		&model.Showtime{},
		&model.User{},
		&model.Reservation{},
	)

		// Create a user, movie and showtime for testing
		user := &model.User{ID: 1, Name: "testuser"}
		db.Create(user)
		movie := &model.Movie{ID: 1, Title: "Test Movie"}
		db.Create(movie)
		showtime := &model.Showtime{ID: 1, MovieID: 1, StartAt: time.Now()}
		db.Create(showtime)
	
		svc := service.NewReservationService(db)
	
		// Test GetRemainingTickets for a new showtime
		remainingTickets, err := svc.GetRemainingTickets(showtime.ID)
		require.NoError(t, err)
		require.Equal(t, model.SeatCount, remainingTickets)
	
		// Test Reserve
		err = svc.Reserve(user.ID, showtime.ID)
		require.NoError(t, err)
	
		// Test GetRemainingTickets after one reservation
		remainingTickets, err = svc.GetRemainingTickets(showtime.ID)
		require.NoError(t, err)
		require.Equal(t, model.SeatCount-1, remainingTickets)
	
		// Test duplicate reservation
		err = svc.Reserve(user.ID, showtime.ID)
		require.ErrorIs(t, err, service.ErrAlreadyReserved)
	
		// Test reservation for non-existent showtime
		err = svc.Reserve(user.ID, 999)
		require.ErrorIs(t, err, service.ErrShowtimeNotExist)
	
		// Test reservation when no tickets are available
		// Fill up the showtime
		for i := 2; i <= model.SeatCount; i++ {
			anotherUser := &model.User{ID: uint(i), Name: fmt.Sprintf("testuser%d", i)}
			db.Create(anotherUser)
			err = svc.Reserve(anotherUser.ID, showtime.ID)
			require.NoError(t, err)
		}
		remainingTickets, err = svc.GetRemainingTickets(showtime.ID)
		require.NoError(t, err)
		require.Equal(t, 0, remainingTickets)
	
		anotherUser := &model.User{ID: 101, Name: "anotheruser"}
	db.Create(anotherUser)
	err = svc.Reserve(anotherUser.ID, showtime.ID)
	require.ErrorIs(t, err, service.ErrNoTicketsAvailable)

	// Test CancelReservation
	// Get a reservation to cancel
	var reservation model.Reservation
	db.Where("user_id = ? AND showtime_id = ?", user.ID, showtime.ID).First(&reservation)
	err = svc.CancelReservation(reservation.ID)
	require.NoError(t, err)

	// Test GetRemainingTickets after cancellation
	remainingTickets, err = svc.GetRemainingTickets(showtime.ID)
	require.NoError(t, err)
	require.Equal(t, 1, remainingTickets)

	// Test canceling a non-existent reservation
	err = svc.CancelReservation(999)
	require.ErrorIs(t, err, service.ErrNotFound)

	// Test GetReservationsByUserID
	reservations, err := svc.GetReservationsByUserID(user.ID)
	require.NoError(t, err)
	// after cancellation, there should be no reservations for this user
	require.Equal(t, 0, len(reservations))

	// re-reserve to test GetReservationByID and GetReservationsByUserID
	err = svc.Reserve(user.ID, showtime.ID)
	require.NoError(t, err)

	reservations, err = svc.GetReservationsByUserID(user.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(reservations))
	require.Equal(t, user.ID, reservations[0].UserID)
	require.Equal(t, showtime.ID, reservations[0].ShowtimeID)

	// Test GetReservationByID
	fetchedReservation, err := svc.GetReservationByID(reservations[0].ID)
	require.NoError(t, err)
	require.Equal(t, reservations[0].ID, fetchedReservation.ID)

	// Test GetReservationByID with non-existent ID
	_, err = svc.GetReservationByID(999)
	require.ErrorIs(t, err, service.ErrNotFound)

	// Test GetReservationsByUserID for a user with no reservations
	anotherUserWithoutReservation := &model.User{ID: 102, Name: "anotheruser2"}
	db.Create(anotherUserWithoutReservation)
	reservations, err = svc.GetReservationsByUserID(anotherUserWithoutReservation.ID)
	require.NoError(t, err)
	require.Equal(t, 0, len(reservations))

	db.Migrator().DropTable(
		&model.Movie{},
		&model.Showtime{},
		&model.User{},
		&model.Reservation{},
	)
}
