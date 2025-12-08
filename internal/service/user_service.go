package service

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/security"
)

type UserService interface {
	CreateUser(userName, password string, role model.UserRole) error
	DeleteUser(userName string, password string) error
	ValidateUser(userName string, password string) (bool, error)
	GetUserRoleByName(userName string) (model.UserRole, error)
	GetUserIDByName(userName string) (uint, error)
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

func (s *userService) CreateUser(userName, password string, role model.UserRole) error {
	_, err := s.repo.GetByName(userName)
	if err == nil {
		return ErrAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}
	return s.repo.Create(&model.User{
		Name:           userName,
		HashedPassword: hash,
		Role:           role,
	})
}

func (s *userService) DeleteUser(userName string, password string) error {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err = s.hasher.Compare(user.HashedPassword, password); err != nil {
		return ErrInvalidCredential
	}
	return s.repo.DeleteByName(userName)
}

func (s *userService) ValidateUser(userName string, password string) (bool, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if err = s.hasher.Compare(user.HashedPassword, password); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *userService) GetUserRoleByName(userName string) (model.UserRole, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return user.Role, nil
}

func (s *userService) GetUserIDByName(userName string) (uint, error) {
	user, err := s.repo.GetByName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return uint(user.ID), nil
}
