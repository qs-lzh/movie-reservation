package service_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/util"

	"github.com/qs-lzh/movie-reservation/internal/service"
)

func TestShowtimeService(t *testing.T) {

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
		model.Movie{},
		model.Showtime{},
	)
	db.AutoMigrate(
		model.Movie{},
		model.Showtime{},
	)

	movieService := service.NewMovieService(db)
	showtimeService := service.NewShowtimeService(db)

	// test CreateShowtime
	var movieID uint = 10002
	movie := &model.Movie{
		ID:    uint(movieID),
		Title: "Harry Potter",
		Description: `Harry Potter is a film series based on the Harry Potter series of novels by J. K. Rowling.
		The series was produced and distributed by Warner Bros.Pictures and consists of eight fantasy films`,
	}
	err = movieService.CreateMovie(movie)
	require.NoError(t, err)
	startTime := time.Now().AddDate(1, 0, 1)
	var hallID uint = 1
	err = showtimeService.CreateShowtime(movieID, startTime, hallID)
	require.NoError(t, err)

	// test GetShowtimesByMovieID and GetShowtimesByMovieID
	showtimes, err := showtimeService.GetShowtimesByMovieID(movieID)
	require.NoError(t, err)
	require.Equal(t, 1, len(showtimes))
	showtimeID := showtimes[0].ID
	showtime, err := showtimeService.GetShowtimeByID(showtimeID)
	require.NoError(t, err)
	require.Equal(t, showtimeID, showtime.ID)

	// test GetAllShowtimes
	allShowtimes, err := showtimeService.GetAllShowtimes()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(allShowtimes), 1)
	// verify the created showtime is in the list
	found := false
	for _, st := range allShowtimes {
		if st.ID == showtimeID {
			found = true
			break
		}
	}
	require.True(t, found, "Created showtime should be in the list of all showtimes")

	// test UpdateShowtime
	hallID = uint(2)
	err = showtimeService.UpdateShowtime(showtimeID, startTime, hallID)
	require.NoError(t, err)
	showtime, err = showtimeService.GetShowtimeByID(showtimeID)
	require.NoError(t, err)
	require.Equal(t, hallID, showtime.HallID)

	// test DeleteShowtimeByID
	err = showtimeService.DeleteShowtimeByID(showtimeID)
	require.NoError(t, err)
	_, err = showtimeService.GetShowtimeByID(showtimeID)
	require.ErrorIs(t, err, service.ErrNotFound)

	db.Migrator().DropTable(
		model.Movie{},
		model.Showtime{},
	)
}
