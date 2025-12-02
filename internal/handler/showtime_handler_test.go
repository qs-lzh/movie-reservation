package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

func TestGetShowtimeByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			showtimes: []model.Showtime{
				{ID: 1, MovieID: 1, StartAt: time.Now(), HallID: 1},
			},
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/showtimes/:id", h.GetShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), `"id":1`)
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/showtimes/:id", h.GetShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/abc", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid showtime id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("not found", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			showtimes: []model.Showtime{{ID: 1, MovieID: 1, StartAt: time.Now(), HallID: 1}},
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/showtimes/:id", h.GetShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/999", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 404, w.Code)
		require.Contains(t, w.Body.String(), "Showtime not exists")
		require.Contains(t, w.Body.String(), "NOT_FOUND")
	})

	t.Run("service error", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			err: errors.New("database error"),
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/showtimes/:id", h.GetShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
	})
}

func TestGetShowtimeAvailability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockReservation := &mockReservationService{
			remainingTickets: 50,
		}
		h := handler.NewShowtimeHandler(&app.App{ReservationService: mockReservation})

		router := gin.New()
		router.GET("/api/showtimes/:id/availability", h.GetShowtimeAvailability)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/1/availability", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), `"remaining_tickets":50`)
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockReservation := &mockReservationService{}
		h := handler.NewShowtimeHandler(&app.App{ReservationService: mockReservation})

		router := gin.New()
		router.GET("/api/showtimes/:id/availability", h.GetShowtimeAvailability)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/abc/availability", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid showtime id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("not found", func(t *testing.T) {
		mockReservation := &mockReservationService{
			err: service.ErrNotFound,
		}
		h := handler.NewShowtimeHandler(&app.App{ReservationService: mockReservation})

		router := gin.New()
		router.GET("/api/showtimes/:id/availability", h.GetShowtimeAvailability)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/999/availability", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 404, w.Code)
		require.Contains(t, w.Body.String(), "Showtime not exists")
		require.Contains(t, w.Body.String(), "NOT_FOUND")
	})

	t.Run("service error", func(t *testing.T) {
		mockReservation := &mockReservationService{
			err: errors.New("database error"),
		}
		h := handler.NewShowtimeHandler(&app.App{ReservationService: mockReservation})

		router := gin.New()
		router.GET("/api/showtimes/:id/availability", h.GetShowtimeAvailability)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/showtimes/1/availability", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
	})
}

// mockServiceForShowtimeTests is a mock implementation of service.ShowtimeService for testing
type mockServiceForShowtimeTests struct {
	showtimes []model.Showtime
	err       error
}

func (m *mockServiceForShowtimeTests) CreateShowtime(movieID uint, startTime time.Time, hallID uint) error {
	return m.err
}
func (m *mockServiceForShowtimeTests) UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error {
	return m.err
}
func (m *mockServiceForShowtimeTests) DeleteShowtimeByID(showtimeID uint) error {
	return m.err
}
func (m *mockServiceForShowtimeTests) GetShowtimeByID(showtimeID uint) (*model.Showtime, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, showtime := range m.showtimes {
		if showtime.ID == showtimeID {
			return &showtime, nil
		}
	}
	return nil, service.ErrNotFound
}
func (m *mockServiceForShowtimeTests) GetShowtimesByMovieID(movieID uint) ([]model.Showtime, error) {
	if m.err != nil {
		return nil, m.err
	}
	var filteredShowtimes []model.Showtime
	for _, st := range m.showtimes {
		if st.MovieID == movieID {
			filteredShowtimes = append(filteredShowtimes, st)
		}
	}
	// If no showtimes were appended, return an empty, non-nil slice
	if filteredShowtimes == nil {
		return []model.Showtime{}, nil
	}
	return filteredShowtimes, nil
}

func TestCreateShowtime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.POST("/showtimes", h.CreateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"movie_id": 1, "start_at": "2024-12-25T14:30:00Z", "hall_id": 1}`
		req, _ := http.NewRequest("POST", "/showtimes", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 201, w.Code)
		require.Contains(t, w.Body.String(), "Showtime created successfully")
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("missing required fields", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.POST("/showtimes", h.CreateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"movie_id": 1, "start_at": "2024-12-25T14:30:00Z"}` // missing hall_id
		req, _ := http.NewRequest("POST", "/showtimes", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid request body")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.POST("/showtimes", h.CreateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{invalid json}`
		req, _ := http.NewRequest("POST", "/showtimes", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid request body")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("service error", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{
			err: errors.New("database error"),
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.POST("/showtimes", h.CreateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"movie_id": 1, "start_at": "2024-12-25T14:30:00Z", "hall_id": 1}`
		req, _ := http.NewRequest("POST", "/showtimes", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
		require.Contains(t, w.Body.String(), `"success":false`)
	})
}

func TestUpdateShowtime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"start_at": "2024-12-26T15:00:00Z", "hall_id": 2}`
		req, _ := http.NewRequest("PUT", "/showtimes/1", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "Showtime updated successfully")
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"start_at": "2024-12-26T15:00:00Z", "hall_id": 2}`
		req, _ := http.NewRequest("PUT", "/showtimes/abc", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid showtime id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("no fields provided for update", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{}` // no fields to update
		req, _ := http.NewRequest("PUT", "/showtimes/1", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "At least one field (start_at or hall_id) must be provided for update")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{invalid json}`
		req, _ := http.NewRequest("PUT", "/showtimes/1", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid request body")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("not found", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{
			err: service.ErrNotFound,
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"start_at": "2024-12-26T15:00:00Z", "hall_id": 2}`
		req, _ := http.NewRequest("PUT", "/showtimes/999", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 404, w.Code)
		require.Contains(t, w.Body.String(), "Showtime not exists")
		require.Contains(t, w.Body.String(), "NOT_FOUND")
	})

	t.Run("service error", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{
			err: errors.New("database error"),
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.PUT("/showtimes/:id", h.UpdateShowtime)

		w := httptest.NewRecorder()
		jsonData := `{"start_at": "2024-12-26T15:00:00Z", "hall_id": 2}`
		req, _ := http.NewRequest("PUT", "/showtimes/1", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
		require.Contains(t, w.Body.String(), `"success":false`)
	})
}

func TestDeleteShowtimeByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.DELETE("/showtimes/:id", h.DeleteShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/showtimes/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "Showtime deleted successfully")
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.DELETE("/showtimes/:id", h.DeleteShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/showtimes/abc", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid showtime id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("not found", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{
			err: service.ErrNotFound,
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.DELETE("/showtimes/:id", h.DeleteShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/showtimes/999", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 404, w.Code)
		require.Contains(t, w.Body.String(), "Showtime not exists")
		require.Contains(t, w.Body.String(), "NOT_FOUND")
	})

	t.Run("service error", func(t *testing.T) {
		mockShowtime := &mockServiceForShowtimeTests{
			err: errors.New("database error"),
		}
		h := handler.NewShowtimeHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.DELETE("/showtimes/:id", h.DeleteShowtimeByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/showtimes/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
		require.Contains(t, w.Body.String(), `"success":false`)
	})
}
