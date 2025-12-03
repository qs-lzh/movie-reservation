package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qs-lzh/movie-reservation/internal/model"
)

func TestReservationFlow(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup(t)

	// Setup: Create admin, users, movie and showtime
	_, adminToken := ts.CreateTestUser(t, "admin", "admin123", "admin")
	_, user1Token := ts.CreateTestUser(t, "user1", "user123", "user")

	movieID := ts.CreateTestMovie(t, adminToken, "Interstellar", "Space movie")

	// Create showtime
	showtimeReq := map[string]interface{}{
		"movie_id": movieID,
		"start_at": "2025-12-15T19:00:00Z",
		"hall_id":  1,
	}
	w := ts.DoRequest("POST", "/showtimes/", showtimeReq, map[string]string{
		"Cookie": "jwt=" + adminToken,
	})
	require.Equal(t, http.StatusCreated, w.Code)

	// Get showtime ID from database
	var showtime model.Showtime
	err := ts.DB.Where("movie_id = ?", movieID).First(&showtime).Error
	require.NoError(t, err)
	showtimeID := showtime.ID

	t.Run("complete reservation flow", func(t *testing.T) {
		// 1. Check showtime availability
		w := ts.DoRequest("GET", fmt.Sprintf("/showtimes/%d/availability", showtimeID), nil, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"remaining_tickets":100`)

		// 2. User1 makes a reservation
		reserveReq := map[string]interface{}{
			"showtime_id": showtimeID,
		}
		w = ts.DoRequest("POST", "/reservations/", reserveReq, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.Equal(t, http.StatusCreated, w.Code)

		// Get reservation ID from database
		var reservation model.Reservation
		err = ts.DB.Where("showtime_id = ?", showtimeID).First(&reservation).Error
		assert.NoError(t, err)
		reservation1ID := reservation.ID
		assert.Greater(t, reservation1ID, uint(0))

		// 3. Check availability decreased
		w = ts.DoRequest("GET", fmt.Sprintf("/showtimes/%d/availability", showtimeID), nil, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"remaining_tickets":99`)

		// 4. User1 checks their reservations
		w = ts.DoRequest("GET", "/reservations/me", nil, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`"id":%d`, reservation1ID))
	})

	t.Run("error cases", func(t *testing.T) {
		// Reserve non-existent showtime
		w := ts.DoRequest("POST", "/reservations/", map[string]interface{}{
			"showtime_id": 99999,
		}, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.NotEqual(t, http.StatusCreated, w.Code)

		// Missing showtime_id
		w = ts.DoRequest("POST", "/reservations/", map[string]interface{}{}, map[string]string{
			"Cookie": "jwt=" + user1Token,
		})
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
