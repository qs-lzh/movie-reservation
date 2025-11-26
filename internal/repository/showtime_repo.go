package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

type ShowtimeRepo interface {
	Create(showtime *model.Showtime) error
	GetByMovieID(movieID uint) ([]*model.Showtime, error)
}

type showtimeRepoGorm struct {
	db *gorm.DB
}

func NewShowtimeRepoGorm(db *gorm.DB) *showtimeRepoGorm {
	return &showtimeRepoGorm{
		db: db,
	}
}

func (r *showtimeRepoGorm) Create(showtime *model.Showtime) error {
	ctx := context.Background()
	if err := gorm.G[model.Showtime](r.db).Create(ctx, showtime); err != nil {
		return err
	}
	return nil
}

func (r *showtimeRepoGorm) GetByMovieID(movieID uint) ([]*model.Showtime, error) {
	ctx := context.Background()
	showtimes, err := gorm.G[*model.Showtime](r.db).Where(&model.Showtime{MovieID: movieID}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return showtimes, nil
}
