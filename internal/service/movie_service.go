package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type MovieService interface {
	CreateMovie(movie *model.Movie) error
	UpdateMovie(movie *model.Movie) error
	DeleteMovieByID(id uint) error
	GetMovieByID(id uint) (*model.Movie, error)
	GetMovieByTitle(title string) (*model.Movie, error)
	GetAllMovies() ([]model.Movie, error)
}

type movieService struct {
	db   *gorm.DB
	repo repository.MovieRepo
}

var _ MovieService = (*movieService)(nil)

func NewMovieService(db *gorm.DB) *movieService {
	return &movieService{
		db:   db,
		repo: repository.NewMovieRepoGorm(db),
	}
}

func (s *movieService) CreateMovie(movie *model.Movie) error {
	if err := s.repo.Create(movie); err != nil {
		return err
	}
	return nil
}

// Update movie by ID
func (s *movieService) UpdateMovie(movie *model.Movie) error {
	// First, verify that the movie with this ID exists
	existingMovie, err := s.repo.GetByID(uint(movie.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// check if the new title is already used by another
	// because the title needs to be unique
	if existingMovie.Title != movie.Title {
		anotherMovie, err := s.repo.GetByTitle(movie.Title)
		if err == nil && anotherMovie != nil && anotherMovie.ID != movie.ID {
			return ErrAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return s.repo.Update(*movie)
}

func (s *movieService) DeleteMovieByID(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *movieService) GetMovieByID(id uint) (*model.Movie, error) {
	movie, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *movieService) GetMovieByTitle(title string) (*model.Movie, error) {
	movie, err := s.repo.GetByTitle(title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *movieService) GetAllMovies() ([]model.Movie, error) {
	movies, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return movies, nil
}
