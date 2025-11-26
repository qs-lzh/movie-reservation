package service

import (
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
)

type MovieService interface {
	CreateMovie(movie *model.Movie) error
	UpdateMovie(movie *model.Movie) error
	DeleteMovieByID(id uint) error
	GetMovieByID(id uint) (*model.Movie, error)
	GetAllMovies() ([]model.Movie, error)
}

type movieService struct {
	db   *gorm.DB
	repo repository.MovieRepo
}

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

func (s *movieService) UpdateMovie(movie *model.Movie) error {
	id := movie.ID
	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}
	if err := s.repo.Create(movie); err != nil {
		return err
	}
	return nil
}

func (s *movieService) DeleteMovieByID(id uint) error {
	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}
	return nil
}

func (s *movieService) GetMovieByID(id uint) (*model.Movie, error) {
	movie, err := s.repo.GetByID(id)
	if err != nil {
		return &model.Movie{}, err
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
