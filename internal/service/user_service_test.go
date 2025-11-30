package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

func TestUserService(t *testing.T) {
	// load test env
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// drop and recreate tables
	db.Migrator().DropTable(
		model.User{},
	)
	db.AutoMigrate(
		model.User{},
	)

	userService := service.NewUserService(db)

	userName := "alice"
	password := "password123"

	// test CreateUser
	err = userService.CreateUser(userName, password)
	require.NoError(t, err)

	// verify user exists in DB
	var user model.User
	result := db.Where(&model.User{Name: userName}).First(&user)
	require.NoError(t, result.Error)
	require.Equal(t, userName, user.Name)
	require.NotEmpty(t, user.HashedPassword)

	// test ValidateUser: correct password
	ok, err := userService.ValidateUser(userName, password)
	require.NoError(t, err)
	require.True(t, ok)

	// test ValidateUser: wrong password
	ok, err = userService.ValidateUser(userName, "wrong-password")
	require.NoError(t, err)
	require.False(t, ok)

	// test ValidateUser: non-existent user
	ok, err = userService.ValidateUser("bob", password)
	require.NoError(t, err)
	require.False(t, ok)

	// test DeleteUser: wrong password should return error
	err = userService.DeleteUser(userName, "wrong-password")
	require.Error(t, err)

	// test DeleteUser: correct password should delete user
	err = userService.DeleteUser(userName, password)
	require.NoError(t, err)
	result = db.Where(&model.User{Name: userName}).First(&user)
	require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)

	// test DeleteUser on non-existent user should be no-op
	err = userService.DeleteUser("non-existent", password)
	require.ErrorIs(t, err, service.ErrNotFound)

	// cleanup
	db.Migrator().DropTable(
		model.User{},
	)
}
