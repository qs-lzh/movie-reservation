import apiClient from './index'

export const hallAPI = {
  getHalls: () => {
    return apiClient.get('/halls/')
  },

  getHallById: (id) => {
    return apiClient.get(`/halls/${id}`)
  },

  createHall: (name, seatCount, rows, cols) => {
    return apiClient.post('/halls/', {
      name,
      seat_count: seatCount,
      rows,
      cols
    })
  },

  updateHall: (id, name, seatCount, rows, cols) => {
    return apiClient.put(`/halls/${id}`, {
      name,
      seat_count: seatCount,
      rows,
      cols
    })
  },

  deleteHall: (id) => {
    return apiClient.delete(`/halls/${id}`)
  }
}