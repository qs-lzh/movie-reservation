package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

// mockUserService implements service.UserService
type mockUserService struct {
	createUserErr  error
	validateUserFn func(username, password string) (bool, error)
	getUserRoleFn  func(username string) (model.UserRole, error)
	getUserIDFn    func(username string) (uint, error)
}

func (m *mockUserService) CreateUser(username, password string, role model.UserRole) error {
	return m.createUserErr
}

func (m *mockUserService) DeleteUser(userName string, password string) error {
	return nil // Not directly tested by AuthHandler, so a no-op is fine
}

func (m *mockUserService) ValidateUser(username, password string) (bool, error) {
	if m.validateUserFn != nil {
		return m.validateUserFn(username, password)
	}
	return false, nil // Default behavior if not set
}

func (m *mockUserService) GetUserRoleByName(userName string) (model.UserRole, error) {
	if m.getUserRoleFn != nil {
		return m.getUserRoleFn(userName)
	}
	return "", nil // Default behavior if not set
}

func (m *mockUserService) GetUserIDByName(userName string) (uint, error) {
	if m.getUserIDFn != nil {
		return m.getUserIDFn(userName)
	}
	return 1, nil // Default behavior if not set
}

// mockAuthService implements service.AuthService
type mockAuthService struct {
	loginToken          string
	loginErr            error
	validateTokenClaims jwt.MapClaims
	validateTokenErr    error
}

func (m *mockAuthService) Login(username, password string) (string, error) {
	return m.loginToken, m.loginErr
}

func (m *mockAuthService) ValidateToken(token string) (jwt.MapClaims, error) {
	return m.validateTokenClaims, m.validateTokenErr
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUserService := &mockUserService{}
		h := handler.NewAuthHandler(&app.App{UserService: mockUserService})

		router := gin.New()
		router.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser", "password": "password", "user_role": "user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Contains(t, w.Body.String(), `"success":true`)
		require.Contains(t, w.Body.String(), `Created user named testuser successfully`)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockUserService := &mockUserService{}
		h := handler.NewAuthHandler(&app.App{UserService: mockUserService})

		router := gin.New()
		router.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser"}` // Missing password
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), `"success":false`)
		require.Contains(t, w.Body.String(), `Invalid request body`)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUserService := &mockUserService{
			createUserErr: service.ErrAlreadyExists,
		}
		h := handler.NewAuthHandler(&app.App{UserService: mockUserService})

		router := gin.New()
		router.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonData := `{"username": "existinguser", "password": "password", "user_role": "user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
		require.Contains(t, w.Body.String(), `"success":false`)
		require.Contains(t, w.Body.String(), `User named existinguser already exists`)
		require.Contains(t, w.Body.String(), `USER_CONFLICTS`)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUserService := &mockUserService{
			createUserErr: errors.New("database error"),
		}
		h := handler.NewAuthHandler(&app.App{UserService: mockUserService})

		router := gin.New()
		router.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser", "password": "password", "user_role": "user"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), `"success":false`)
		require.Contains(t, w.Body.String(), `Failed to create user`)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthService := &mockAuthService{
			loginToken: "mock_jwt_token",
		}
		mockUserService := &mockUserService{
			getUserRoleFn: func(username string) (model.UserRole, error) {
				return model.RoleUser, nil
			},
		}
		h := handler.NewAuthHandler(&app.App{
			AuthService: mockAuthService,
			UserService: mockUserService,
		})

		router := gin.New()
		router.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser", "password": "password"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"success":true`)
		require.Contains(t, w.Body.String(), `Login successfully`)

		// Check for the cookie
		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		require.Equal(t, "jwt", cookies[0].Name)
		require.Equal(t, "mock_jwt_token", cookies[0].Value)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockAuthService := &mockAuthService{}
		h := handler.NewAuthHandler(&app.App{AuthService: mockAuthService})

		router := gin.New()
		router.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser"}` // Missing password
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), `"success":false`)
		require.Contains(t, w.Body.String(), `Invalid request body`)
	})

	t.Run("internal server error during login", func(t *testing.T) {
		mockAuthService := &mockAuthService{
			loginErr: errors.New("authentication failed"),
		}
		h := handler.NewAuthHandler(&app.App{AuthService: mockAuthService})

		router := gin.New()
		router.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonData := `{"username": "testuser", "password": "password"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), `"success":false`)
		require.Contains(t, w.Body.String(), `Failed to login`)
	})
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		h := handler.NewAuthHandler(&app.App{}) // AuthService not needed for logout

		router := gin.New()
		router.POST("/logout", h.Logout)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logout", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"success":true`)
		require.Contains(t, w.Body.String(), `Logged out successfully`)

		// Check for the expired cookie
		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		require.Equal(t, "jwt", cookies[0].Name)
		require.Equal(t, "", cookies[0].Value)
		require.Equal(t, -1, cookies[0].MaxAge)
	})
}
