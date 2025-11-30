package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/repository"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func TestUserRepo(t *testing.T) {
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// drop and recreate table
	db.Migrator().DropTable(
		model.User{},
	)
	db.AutoMigrate(
		model.User{},
	)

	userRepo := repository.NewUserRepoGorm(db)

	// test Create and GetByName
	userName := "testuser"
	user := &model.User{
		Name:           userName,
		HashedPassword: "hashed_password_123",
	}

	err = userRepo.Create(user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	// verify user is properly stored
	retrievedUser, err := userRepo.GetByName(userName)
	require.NoError(t, err)
	require.Equal(t, userName, retrievedUser.Name)
	require.Equal(t, "hashed_password_123", retrievedUser.HashedPassword)
	require.Equal(t, user.ID, retrievedUser.ID)

	// test GetByName for non-existent user
	_, err = userRepo.GetByName("nonexistent_user")
	require.Error(t, err)

	// test Create with duplicate name (should work for this test since we're testing repo methods)
	anotherUser := &model.User{
		Name:           userName, // same name as before
		HashedPassword: "different_hashed_password",
	}
	err = userRepo.Create(anotherUser)
	require.Error(t, err) // should error due to unique constraint

	// test DeleteByName
	err = userRepo.DeleteByName(userName)
	require.NoError(t, err)

	// verify user is deleted
	_, err = userRepo.GetByName(userName)
	require.Error(t, err)

	// test DeleteByName for non-existent user (should not error)
	err = userRepo.DeleteByName("nonexistent_user")
	require.NoError(t, err)

	// drop all tables
	db.Migrator().DropTable(
		model.User{},
	)
}