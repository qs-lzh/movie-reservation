package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieFlow(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup(t)

	// Create admin and regular user
	_, adminToken := ts.CreateTestUser(t, "admin", "admin123", "admin")
	_, userToken := ts.CreateTestUser(t, "user1", "user123", "user")

	t.Run("complete movie CRUD flow", func(t *testing.T) {
		// 1. Create movie (admin only)
		movieReq := map[string]interface{}{
			"title":       "The Matrix",
			"description": "A great sci-fi movie",
		}
		w := ts.DoRequest("POST", "/movies/", movieReq, map[string]string{
			"Cookie": "jwt=" + adminToken,
		})
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResp struct {
			Data struct {
				ID          uint   `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &createResp)
		movieID := createResp.Data.ID
		assert.Greater(t, movieID, uint(0))

		// 2. Regular user cannot create movie
		w = ts.DoRequest("POST", "/movies/", movieReq, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusForbidden, w.Code)

		// 3. Get movie by ID (any authenticated user)
		w = ts.DoRequest("GET", fmt.Sprintf("/movies/%d", movieID), nil, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "The Matrix")

		// 4. Get all movies
		w = ts.DoRequest("GET", "/movies/", nil, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "The Matrix")

		// 5. Delete movie (admin only)
		w = ts.DoRequest("DELETE", fmt.Sprintf("/movies/%d", movieID), nil, map[string]string{
			"Cookie": "jwt=" + adminToken,
		})
		assert.Equal(t, http.StatusOK, w.Code)

		// 6. Verify movie is deleted
		w = ts.DoRequest("GET", fmt.Sprintf("/movies/%d", movieID), nil, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("error cases", func(t *testing.T) {
		// Invalid movie ID
		w := ts.DoRequest("GET", "/movies/abc", nil, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Movie not found
		w = ts.DoRequest("GET", "/movies/99999", nil, map[string]string{
			"Cookie": "jwt=" + userToken,
		})
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Create movie without title
		w = ts.DoRequest("POST", "/movies/", map[string]interface{}{
			"description": "No title",
		}, map[string]string{
			"Cookie": "jwt=" + adminToken,
		})
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMovieShowtimes(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup(t)

	_, adminToken := ts.CreateTestUser(t, "admin", "admin123", "admin")
	_, userToken := ts.CreateTestUser(t, "user1", "user123", "user")

	// Create a movie
	movieID := ts.CreateTestMovie(t, adminToken, "Inception", "Dream heist movie")

	// Create showtimes for the movie
	showtimeReq := map[string]interface{}{
		"movie_id":  movieID,
		"start_at":  "2025-12-10T20:00:00Z",
		"hall_id":   1,
	}
	w := ts.DoRequest("POST", "/showtimes/", showtimeReq, map[string]string{
		"Cookie": "jwt=" + adminToken,
	})
	require.Equal(t, http.StatusCreated, w.Code)

	// Get movie showtimes
	w = ts.DoRequest("GET", fmt.Sprintf("/movies/%d/showtimes", movieID), nil, map[string]string{
		"Cookie": "jwt=" + userToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), fmt.Sprintf(`"movie_id":%d`, movieID))
}
