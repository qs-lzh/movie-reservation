package repository

import "gorm.io/gorm"

type MovieRepo interface {
	Create(movie *Movie) error
	GetByID(id int64) (*Movie, error)
}

type movieRepoGorm struct {
	db *gorm.DB
}
