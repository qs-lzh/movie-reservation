package repository

import (
	"context"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user *model.User) error
	GetHashedPassword(name string) (string, error)
}

type userRepoGorm struct {
	db *gorm.DB
}

var _ UserRepo = (*userRepoGorm)(nil)

func NewUserRepoGorm(db *gorm.DB) *userRepoGorm {
	return &userRepoGorm{
		db: db,
	}
}

func (r *userRepoGorm) Create(user *model.User) error {
	ctx := context.Background()
	if err := gorm.G[model.User](r.db).Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *userRepoGorm) GetHashedPassword(name string) (string, error) {
	ctx := context.Background()
	user, err := gorm.G[model.User](r.db).Where(model.User{Name: name}).First(ctx)
	if err != nil {
		return "", nil
	}
	hashedPassword := user.HashedPassword
	return hashedPassword, nil
}
