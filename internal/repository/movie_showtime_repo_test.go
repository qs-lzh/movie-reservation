package repository_test

/*
* Test Movie and Showtime together
 */

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func TestMovieShowtimeRepo(t *testing.T) {
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("DATABASE_DSN")

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

	movieRepo := repository.NewMovieRepoGorm(db)
	showtimeRepo := repository.NewShowtimeRepoGorm(db)

	// create a movie and several showtimes belong to it, store them in database
	// and validate it is properly stored

	var movieID uint = 10002
	// store movie without showtimes
	movie := &model.Movie{
		ID:    movieID,
		Title: "Harry Potter",
		Description: `Harry Potter is a film series based on the Harry Potter series of novels by J. K. Rowling.
		The series was produced and distributed by Warner Bros.Pictures and consists of eight fantasy films`,
	}

	err = movieRepo.Create(movie)
	require.NoError(t, err)

	var showtimes []model.Showtime
	for i := range 5 {
		showtime := model.Showtime{
			MovieID: movieID,
			StartAt: time.Now().AddDate(1, 0, i),
		}
		err := showtimeRepo.Create(showtime)
		require.NoError(t, err)
		showtimes = append(showtimes, showtime)
	}

	movie, err = movieRepo.GetByID(movieID)
	require.NoError(t, err)
	require.Equal(t, movieID, movie.ID)
	require.Equal(t, "Harry Potter", movie.Title)

	retrievedShowtimes, err := showtimeRepo.GetByMovieID(movieID)
	require.NoError(t, err)
	require.Equal(t, 5, len(retrievedShowtimes))

	for _, showtime := range retrievedShowtimes {
		require.Equal(t, movieID, showtime.MovieID)
	}

	// test delete

	err = showtimeRepo.DeleteByID(1)
	require.NoError(t, err)
	// validate delete showtime properly
	showtime, err := showtimeRepo.GetByMovieID(1)
	require.NoError(t, err)
	require.Equal(t, 0, len(showtime))
	//delete movie
	err = showtimeRepo.DeleteByMovieID(movieID)
	require.NoError(t, err)
	err = movieRepo.DeleteByID(movieID)
	require.NoError(t, err)

	// drop all tables
	db.Migrator().DropTable(
		model.Movie{},
		model.Showtime{},
	)
}
