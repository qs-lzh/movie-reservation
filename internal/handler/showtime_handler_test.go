package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
