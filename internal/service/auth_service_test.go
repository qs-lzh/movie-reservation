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

func TestAuthService(t *testing.T) {
	// load test env
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	})
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

		authService := service.NewJWTAuthService(userService)

	

		userName := "alice"

		password := "password123"

	

		// create a user

		err = userService.CreateUser(userName, password)

		require.NoError(t, err)

	

		// test Login

		token, err := authService.Login(userName, password)

		require.NoError(t, err)

		require.NotEmpty(t, token)

	

		// test Login with wrong password

		_, err = authService.Login(userName, "wrong-password")

		require.NoError(t, err)

	

		// test ValidateToken with correct token

		err = authService.ValidateToken(token)

		require.NoError(t, err)

	

		// test Authenticate with wrong token

		err = authService.ValidateToken("wrong-token")

		require.Error(t, err)

	

		// cleanup

		db.Migrator().DropTable(

			model.User{},

		)

	}
