package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

type MovieRepo interface {
	Create(movie *model.Movie) error
	GetByID(id uint) (*model.Movie, error)
}

type movieRepoGorm struct {
	db *gorm.DB
}

func NewMovieRepoGorm(db *gorm.DB) *movieRepoGorm {
	return &movieRepoGorm{
		db: db,
	}
}

func (r *movieRepoGorm) Create(movie *model.Movie) error {
	ctx := context.Background()
	if err := gorm.G[model.Movie](r.db).Create(ctx, movie); err != nil {
		return err
	}
	return nil
}

func (r *movieRepoGorm) GetByID(id uint) (*model.Movie, error) {
	ctx := context.Background()
	movie, err := gorm.G[model.Movie](r.db).Where(&model.Movie{ID: id}).First(ctx)
	if err != nil {
		return &model.Movie{}, err
	}
	return &movie, nil
}
