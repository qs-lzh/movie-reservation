import apiClient from './index'

export const movieAPI = {
  getMovies: () => {
    return apiClient.get('/movies/')
  },

  getMovieById: (id) => {
    return apiClient.get(`/movies/${id}`)
  },

  getMovieShowtimes: (id) => {
    return apiClient.get(`/movies/${id}/showtimes`)
  },

  createMovie: (title, description) => {
    return apiClient.post('/movies/', {
      title,
      description
    })
  },

  updateMovie: (id, title, description) => {
    return apiClient.put(`/movies/${id}`, {
      title,
      description
    })
  },

  deleteMovie: (id) => {
    return apiClient.delete(`/movies/${id}`)
  }
}
