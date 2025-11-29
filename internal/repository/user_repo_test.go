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

	db.Migrator().DropTable(
		model.User{},
	)
	db.AutoMigrate(
		model.User{},
	)

	repo := repository.NewUserRepoGorm(db)

	user := &model.User{
		Name:           "alice",
		HashedPassword: "hashed_pw_123",
	}
	// Create user
	err = repo.Create(user)
	require.NoError(t, err)

	// GetHashedPassword: existing user
	hp, err := repo.GetHashedPassword("alice")
	require.NoError(t, err)
	require.Equal(t, "hashed_pw_123", hp)

	// GetHashedPassword: non-existent user
	hp, err = repo.GetHashedPassword("bob")
	require.NoError(t, err)
	require.Equal(t, "", hp)

	// GetByName: existing user
	gotUser, err := repo.GetByName("alice")
	require.NoError(t, err)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.HashedPassword, gotUser.HashedPassword)

	// Delete: existing user
	err = repo.Delete("alice")
	require.NoError(t, err)

	// After delete, GetHashedPassword should return empty string
	hp, err = repo.GetHashedPassword("alice")
	require.NoError(t, err)
	require.Equal(t, "", hp)

	// After delete, GetByName should return error (record not found)
	_, err = repo.GetByName("alice")
	require.Error(t, err)

	// Delete: non-existent user should be no-op
	err = repo.Delete("non-existent")
	require.NoError(t, err)

	db.Migrator().DropTable(
		model.User{},
	)
}
