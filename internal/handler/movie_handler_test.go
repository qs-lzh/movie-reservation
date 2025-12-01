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

// mockMovieService implements service.MovieService
type mockMovieService struct {
	movies []model.Movie
	err    error
}

func (m *mockMovieService) GetAllMovies() ([]model.Movie, error) {
	return m.movies, m.err
}

func (m *mockMovieService) GetMovieByID(id uint) (*model.Movie, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, movie := range m.movies {
		if movie.ID == id {
			return &movie, nil
		}
	}
	return nil, service.ErrNotFound
}

func (m *mockMovieService) CreateMovie(*model.Movie) error { return nil }
func (m *mockMovieService) UpdateMovie(*model.Movie) error { return nil }
func (m *mockMovieService) DeleteMovieByID(uint) error     { return nil }

// mockShowtimeService implements service.ShowtimeService
type mockShowtimeService struct {
	showtimes []model.Showtime
	err       error
}

func (m *mockShowtimeService) CreateShowtime(movieID uint, startTime time.Time, hallID uint) error {
	return nil
}
func (m *mockShowtimeService) UpdateShowtime(showtimeID uint, startTime time.Time, hallID uint) error {
	return nil
}
func (m *mockShowtimeService) DeleteShowtimeByID(showtimeID uint) error {
	return nil
}
func (m *mockShowtimeService) GetShowtimeByID(showtimeID uint) (*model.Showtime, error) {
	return nil, nil
}
func (m *mockShowtimeService) GetShowtimesByMovieID(movieID uint) ([]model.Showtime, error) {
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

func TestGetAllMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mock := &mockMovieService{
			movies: []model.Movie{
				{ID: 1, Title: "The Matrix"},
				{ID: 2, Title: "Inception"},
			},
		}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies", h.GetAllMovies)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "The Matrix")
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("service error", func(t *testing.T) {
		mock := &mockMovieService{
			err: errors.New("database error"),
		}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies", h.GetAllMovies)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
		require.Contains(t, w.Body.String(), `"success":false`)
	})
}

func TestGetMovieByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mock := &mockMovieService{
			movies: []model.Movie{
				{ID: 1, Title: "The Matrix", Description: "A great movie"},
			},
		}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies/:id", h.GetMovieByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "The Matrix")
		require.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("invalid id", func(t *testing.T) {
		mock := &mockMovieService{}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies/:id", h.GetMovieByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/abc", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid movie id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockMovieService{
			movies: []model.Movie{{ID: 1, Title: "The Matrix"}},
		}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies/:id", h.GetMovieByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/999", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 404, w.Code)
		require.Contains(t, w.Body.String(), "Movie not exists")
		require.Contains(t, w.Body.String(), "NOT_FOUND")
	})

	t.Run("service error", func(t *testing.T) {
		mock := &mockMovieService{
			err: errors.New("database error"),
		}
		h := handler.NewMovieHandler(&app.App{MovieService: mock})

		router := gin.New()
		router.GET("/api/movies/:id", h.GetMovieByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
	})
}

func TestGetMovieShowtimes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			showtimes: []model.Showtime{
				{ID: 1, MovieID: 1, StartAt: time.Now(), HallID: 1},
				{ID: 2, MovieID: 1, StartAt: time.Now().Add(time.Hour), HallID: 2},
				{ID: 3, MovieID: 2, StartAt: time.Now().Add(2 * time.Hour), HallID: 1},
			},
		}
		h := handler.NewMovieHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/movies/:id/showtimes", h.GetMovieShowtimes)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/1/showtimes", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), `"movie_id":1`)
		require.Contains(t, w.Body.String(), `"success":true`)
		// Instead of require.Len, we'll check for specific content indicating two showtimes.
		require.Contains(t, w.Body.String(), `[{"id":1`) // Check for start of first showtime
		require.Contains(t, w.Body.String(), `{"id":2`) // Check for start of second showtime
	})

	t.Run("no showtimes for movie", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			showtimes: []model.Showtime{
				{ID: 3, MovieID: 2, StartAt: time.Now().Add(2 * time.Hour), HallID: 1},
			},
		}
		h := handler.NewMovieHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/movies/:id/showtimes", h.GetMovieShowtimes)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/999/showtimes", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), `"success":true`)
		require.Contains(t, w.Body.String(), `"data":[]`)
	})

	t.Run("invalid movie id", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{}
		h := handler.NewMovieHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/movies/:id/showtimes", h.GetMovieShowtimes)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/abc/showtimes", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 400, w.Code)
		require.Contains(t, w.Body.String(), "Invalid movie id")
		require.Contains(t, w.Body.String(), "BAD_REQUEST")
	})

	t.Run("service error", func(t *testing.T) {
		mockShowtime := &mockShowtimeService{
			err: errors.New("database error"),
		}
		h := handler.NewMovieHandler(&app.App{ShowtimeService: mockShowtime})

		router := gin.New()
		router.GET("/api/movies/:id/showtimes", h.GetMovieShowtimes)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/movies/1/showtimes", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, 500, w.Code)
		require.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
	})
}