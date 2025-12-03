import apiClient from './index'

export const reservationAPI = {
  createReservation: (showtimeId) => {
    // 确保showtimeId是数字类型
    return apiClient.post('/reservations/', {
      showtime_id: parseInt(showtimeId, 10)
    })
  },

  getUserReservations: () => {
    return apiClient.get('/reservations/me')
  },

  cancelReservation: (reservationId) => {
    return apiClient.delete(`/reservations/${reservationId}`)
  }
}