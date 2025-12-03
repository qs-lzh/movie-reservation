package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/qs-lzh/movie-reservation/config"
	"github.com/qs-lzh/movie-reservation/interfaces/web"
	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/security"
	"github.com/qs-lzh/movie-reservation/internal/util"
)

type TestServer struct {
	App    *app.App
	Router *gin.Engine
	DB     *gorm.DB
}

// SetupTestServer creates a new test server with test PostgreSQL database
func SetupTestServer(t *testing.T) *TestServer {
	gin.SetMode(gin.TestMode)

	// Load environment variables and get test database DSN
	err := util.LoadEnv()
	require.NoError(t, err)
	dsn := os.Getenv("TEST_DATABASE_DSN")
	require.NotEmpty(t, dsn, "TEST_DATABASE_DSN environment variable must be set")

	// Use PostgreSQL test database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	require.NoError(t, err)

	// Drop and recreate tables for clean state
	db.Migrator().DropTable(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
	)

	err = db.AutoMigrate(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
	)
	require.NoError(t, err)

	// Initialize JWT for testing
	cfg := &config.Config{
		JWTSecretKey: "test-secret-key-for-integration-tests",
	}
	security.InitJWT(cfg.JWTSecretKey)

	// Create app instance
	testApp := app.New(cfg, db)

	// Initialize router
	router := web.InitRouter(testApp)

	return &TestServer{
		App:    testApp,
		Router: router,
		DB:     db,
	}
}

// Cleanup drops tables and closes the database connection
func (ts *TestServer) Cleanup(t *testing.T) {
	// Drop all tables
	ts.DB.Migrator().DropTable(
		&model.User{},
		&model.Movie{},
		&model.Showtime{},
		&model.Reservation{},
	)

	sqlDB, err := ts.DB.DB()
	require.NoError(t, err)
	sqlDB.Close()
}

// DoRequest performs an HTTP request and returns the response
func (ts *TestServer) DoRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		if k == "Cookie" {
			req.Header.Set("Cookie", v)
		} else {
			req.Header.Set(k, v)
		}
	}

	w := httptest.NewRecorder()
	ts.Router.ServeHTTP(w, req)
	return w
}

// CreateTestUser creates a test user and returns the user and token
func (ts *TestServer) CreateTestUser(t *testing.T, name, password string, role model.UserRole) (uint, string) {
	registerReq := map[string]interface{}{
		"username":  name,
		"password":  password,
		"user_role": role,
	}

	w := ts.DoRequest("POST", "/users/register", registerReq, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	// Get user ID from database
	var user model.User
	err := ts.DB.Where("name = ?", name).First(&user).Error
	require.NoError(t, err)

	// Login to get token
	loginReq := map[string]interface{}{
		"username": name,
		"password": password,
	}

	w = ts.DoRequest("POST", "/users/login", loginReq, nil)
	require.Equal(t, http.StatusOK, w.Code)

	// Extract token from cookie
	cookies := w.Result().Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			token = cookie.Value
			break
		}
	}
	require.NotEmpty(t, token, "JWT token not found in cookies")

	return user.ID, token
}

// CreateTestMovie creates a test movie and returns the movie ID
func (ts *TestServer) CreateTestMovie(t *testing.T, token, title, description string) uint {
	movieReq := map[string]interface{}{
		"title":       title,
		"description": description,
	}

	headers := map[string]string{
		"Cookie": "jwt=" + token,
	}

	w := ts.DoRequest("POST", "/movies/", movieReq, headers)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	return resp.Data.ID
}
