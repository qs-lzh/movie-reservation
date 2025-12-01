package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
	"github.com/qs-lzh/movie-reservation/internal/model"
	"github.com/qs-lzh/movie-reservation/internal/service"
)

type mockReservationService struct {
	service.ReservationService
	reservations []model.Reservation
	err          error
	remainingTickets int // Added this field
}

func (m *mockReservationService) Reserve(userID uint, showtimeID uint) error {
	return m.err
}

func (m *mockReservationService) GetReservationsByUserID(userID uint) ([]model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	var userReservations []model.Reservation
	for _, r := range m.reservations {
		if r.UserID == userID {
			userReservations = append(userReservations, r)
		}
	}
	return userReservations, nil
}

func (m *mockReservationService) GetReservationByID(reservationID uint) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, r := range m.reservations {
		if r.ID == reservationID {
			return &r, nil
		}
	}
	return nil, service.ErrNotFound
}

func (m *mockReservationService) CancelReservation(reservationID uint) error {
	if m.err != nil {
		return m.err
	}
	for i, r := range m.reservations {
		if r.ID == reservationID {
			m.reservations = append(m.reservations[:i], m.reservations[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound // Simulate not found if not in mock list
}

func (m *mockReservationService) GetRemainingTickets(showtimeID uint) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.remainingTickets, nil
}

func setupReservationRouter(mockService service.ReservationService) *gin.Engine {
	r := gin.Default()
	h := handler.NewReservationHandler(&app.App{ReservationService: mockService})
	group := r.Group("/reservations")
	{
		group.POST("/", h.CreateReservation)
		group.GET("/me", h.GetMyReservations)
		group.DELETE("/:id", h.CancelReservation)
	}
	return r
}

func TestCreateReservation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockReservationService{}
	router := setupReservationRouter(mockService)

	body, _ := json.Marshal(gin.H{"showtime_id": 1})
	req, _ := http.NewRequest("POST", "/reservations/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), "Reservation created successfully")
}

func TestGetMyReservations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockReservationService{
		reservations: []model.Reservation{
			{ID: 1, UserID: 1, ShowtimeID: 101},
			{ID: 2, UserID: 2, ShowtimeID: 102},
			{ID: 3, UserID: 1, ShowtimeID: 103},
		},
	}
	router := setupReservationRouter(mockService)

	req, _ := http.NewRequest("GET", "/reservations/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":1`)
	require.Contains(t, w.Body.String(), `"id":3`)
	require.NotContains(t, w.Body.String(), `"id":2`)
}

func TestCancelReservation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockReservationService{
			reservations: []model.Reservation{{ID: 1, UserID: 1}},
		}
		router := setupReservationRouter(mockService)
		req, _ := http.NewRequest("DELETE", "/reservations/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "Reservation cancelled successfully")
	})

	t.Run("forbidden", func(t *testing.T) {
		mockService := &mockReservationService{
			reservations: []model.Reservation{{ID: 1, UserID: 2}}, // Belongs to another user
		}
		router := setupReservationRouter(mockService)
		req, _ := http.NewRequest("DELETE", "/reservations/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := &mockReservationService{} // No reservations
		router := setupReservationRouter(mockService)
		req, _ := http.NewRequest("DELETE", "/reservations/99", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}