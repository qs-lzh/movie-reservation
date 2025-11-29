package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func TestReservationRepo(t *testing.T) {
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	db.Migrator().DropTable(
		model.Reservation{},
	)
	db.AutoMigrate(
		model.Reservation{},
	)

	repo := repository.NewReservationRepoGorm(db)

	userID := uint(20001)
	showtimeID := uint(30001)
	reservation1 := &model.Reservation{
		UserID:     userID,
		ShowtimeID: showtimeID,
	}
	reservation2 := &model.Reservation{
		UserID:     userID,
		ShowtimeID: showtimeID,
	}

	err = repo.Create(reservation1)
	require.NoError(t, err)
	err = repo.Create(reservation2)
	require.NoError(t, err)

	got, err := repo.GetByID(reservation1.ID)
	require.NoError(t, err)
	require.Equal(t, reservation1.ID, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, showtimeID, got.ShowtimeID)

	byUser, err := repo.GetByUserID(userID)
	require.NoError(t, err)
	require.Equal(t, 2, len(byUser))
	for _, r := range byUser {
		require.Equal(t, userID, r.UserID)
	}

	byShowtime, err := repo.GetByShowtimeID(showtimeID)
	require.NoError(t, err)
	require.Equal(t, 2, len(byShowtime))
	for _, r := range byShowtime {
		require.Equal(t, showtimeID, r.ShowtimeID)
	}

	err = repo.DeleteByID(reservation1.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(reservation1.ID)
	require.Error(t, err)

	db.Migrator().DropTable(
		model.Reservation{},
	)
}

