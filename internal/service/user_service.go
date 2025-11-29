package service

import (
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/security"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(userName string, password string) error
	DeleteUser(userName string, password string) error
	ValidateUser(userName string, password string) (bool, error)
}

type userService struct {
	db     *gorm.DB
	repo   repository.UserRepo
	hasher security.PasswordHasher
}

var _ UserService = (*userService)(nil)

func NewUserService(db *gorm.DB) *userService {
	return &userService{
		db:     db,
		repo:   repository.NewUserRepoGorm(db),
		hasher: security.NewBcryptHasher(10),
	}
}

func (s *userService) CreateUser(userName string, password string) error {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}
	user := &model.User{
		Name:           userName,
		HashedPassword: hash,
	}
	return s.repo.Create(user)
}

func (s *userService) DeleteUser(userName string, password string) error {
	// check if user exists
	user, err := s.repo.GetByName(userName)
	if err != nil {
		// the returned will be gorm.ErrRecordNotFound
		return err
	}
	if user == nil {
		return nil
	}

	hashedPassword, err := s.repo.GetHashedPassword(userName)
	if err != nil {
		return err
	}
	if err := s.hasher.Compare(hashedPassword, password); err != nil {
		return err
	}
	return s.repo.Delete(userName)
}

func (s *userService) ValidateUser(userName string, password string) (bool, error) {
	hashedPassword, err := s.repo.GetHashedPassword(userName)
	if err != nil {
		return false, err
	}
	if hashedPassword == "" {
		return false, nil
	}
	if err := s.hasher.Compare(hashedPassword, password); err != nil {
		return false, nil
	}
	return true, nil
}
