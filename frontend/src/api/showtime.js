import apiClient from './index'

export const showtimeAPI = {
  getShowtimes: () => {
    return apiClient.get('/showtimes/')
  },

  getShowtimeById: (id) => {
    return apiClient.get(`/showtimes/${id}`)
  },

  getShowtimeAvailability: (id) => {
    return apiClient.get(`/showtimes/${id}/availability`)
  },

  createShowtime: (movieId, startAt, hallId) => {
    return apiClient.post('/showtimes/', {
      movie_id: movieId,
      start_at: startAt,
      hall_id: hallId
    })
  },

  updateShowtime: (id, startAt, hallId) => {
    return apiClient.put(`/showtimes/${id}`, {
      start_at: startAt,
      hall_id: hallId
    })
  },

  deleteShowtime: (id) => {
    return apiClient.delete(`/showtimes/${id}`)
  }
}