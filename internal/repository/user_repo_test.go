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
	err = repo.Create(user)
	require.NoError(t, err)
	hp, err := repo.GetHashedPassword("alice")
	require.NoError(t, err)
	require.Equal(t, "hashed_pw_123", hp)
	hp, err = repo.GetHashedPassword("bob")
	require.NoError(t, err)
	require.Equal(t, "", hp)

	db.Migrator().DropTable(
		model.User{},
	)
}
