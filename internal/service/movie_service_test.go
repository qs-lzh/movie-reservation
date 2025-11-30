package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func TestMovieService(t *testing.T) {
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

	svc := service.NewMovieService(db)

	// test CreateMovie
	var movieID uint = 10002
	movie := &model.Movie{
		ID:    uint(movieID),
		Title: "Harry Potter",
		Description: `Harry Potter is a film series based on the Harry Potter series of novels by J. K. Rowling.
		The series was produced and distributed by Warner Bros.Pictures and consists of eight fantasy films`,
	}
	err = svc.CreateMovie(movie)
	require.NoError(t, err)

	// test GetMovieByID
	movie, err = svc.GetMovieByID(movieID)
	require.NoError(t, err)
	require.Equal(t, movieID, movie.ID)

	// test UpdateMovie
	movieID = uint(10003)
	originalMovie := &model.Movie{
		ID:    movieID,
		Title: "Original Movie",
	}
	err = svc.CreateMovie(originalMovie)
	require.NoError(t, err)

	updatedMovie := &model.Movie{
		ID:    movieID,
		Title: "Original Movie", // Use same title to find the movie
		Description: "Updated description",
	}
	err = svc.UpdateMovie(updatedMovie)
	require.NoError(t, err)
	movie, err = svc.GetMovieByID(movieID)
	require.NoError(t, err)
	require.Equal(t, movieID, movie.ID)
	require.Equal(t, "Updated description", movie.Description)

	// test DeleteMovieByID
	err = svc.DeleteMovieByID(movieID)
	require.NoError(t, err)
	movie, err = svc.GetMovieByID(movieID)
	require.ErrorIs(t, err, service.ErrNotFound)

	// test GetAllMovies
	err = svc.CreateMovie(&model.Movie{ID: uint(10005)})
	require.NoError(t, err)
	movies, err := svc.GetAllMovies()
	require.NoError(t, err)
	require.Equal(t, 2, len(movies))

	db.Migrator().DropTable(
		model.Movie{},
		model.Showtime{},
	)
}
