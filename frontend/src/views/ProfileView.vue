<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <div class="profile-header">
        <div class="user-avatar">
          <el-icon :size="60">
            <User />
          </el-icon>
        </div>
        <div class="user-info">
          <h2>Welcome, {{ userStore.user?.username || 'User' }}!</h2>
          <p>User Role: <el-tag :type="getRoleTagType(userStore.user?.role)">{{ userStore.user?.role || 'user' }}</el-tag></p>
          <p>Member since {{ formatDate(new Date()) }}</p>
        </div>
      </div>
    </el-card>

    <el-card class="reservations-card">
      <h3>Your Reservations</h3>

      <div v-if="isLoading" class="loading-container">
        <el-skeleton :rows="4" animated />
      </div>

      <el-empty
        v-else-if="reservations.length === 0"
        description="No reservations yet"
      />

      <div v-else class="reservations-list">
        <el-table
          :data="reservations"
          style="width: 100%"
          stripe
        >
          <el-table-column prop="id" label="Reservation ID" width="120" />
          <el-table-column label="Movie" width="200">
            <template #default="{ row }">
              {{ getMovieTitle(row.showtime) }}
            </template>
          </el-table-column>
          <el-table-column label="Date & Time" width="200">
            <template #default="{ row }">
              {{ formatDate(row.showtime.start_at) }}<br>
              {{ formatTime(row.showtime.start_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="showtime.hall_id" label="Hall" width="100" />
          <el-table-column label="Status" width="120">
            <template #default>
              <el-tag type="success">Confirmed</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="150">
            <template #default="{ row }">
              <el-button
                size="small"
                type="danger"
                @click="cancelReservation(row.id)"
                :disabled="isCancelling"
              >
                Cancel
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { reservationAPI } from '@/api/reservation'
import { showtimeAPI } from '@/api/showtime'
import { movieAPI } from '@/api/movie'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'

const userStore = useUserStore()
const reservations = ref([])
const isLoading = ref(true)
const isCancelling = ref(false)

// 获取用户预订列表
const fetchReservations = async () => {
  try {
    const response = await reservationAPI.getUserReservations()
    reservations.value = response.data.data || []

    // 获取每个预订对应的场次信息和电影信息
    for (const reservation of reservations.value) {
      try {
        // 获取场次信息
        const showtimeResponse = await showtimeAPI.getShowtimeById(reservation.showtime_id)
        reservation.showtime = showtimeResponse.data.data

        // 获取电影信息
        if (reservation.showtime?.movie_id) {
          const movieResponse = await movieAPI.getMovieById(reservation.showtime.movie_id)
          reservation.showtime.movie = movieResponse.data.data
        }
      } catch (error) {
        console.error(`Failed to fetch details for reservation ${reservation.id}:`, error)
      }
    }
  } catch (error) {
    console.error('Failed to fetch reservations:', error)
    let message = 'Failed to load reservations'
    if (error.response?.data?.error?.message) {
      message = error.response.data.error.message
    } else if (error.response?.data?.message) {
      message = error.response.data.message
    }
    ElMessage.error(message)
  } finally {
    isLoading.value = false
  }
}

// 取消预订
const cancelReservation = async (reservationId) => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to cancel this reservation?',
      'Confirm Cancel',
      {
        confirmButtonText: 'Yes',
        cancelButtonText: 'No',
        type: 'warning',
      }
    )

    isCancelling.value = true

    await reservationAPI.cancelReservation(reservationId)
    ElMessage.success('Reservation cancelled successfully')

    // 刷新预订列表
    await fetchReservations()
  } catch (error) {
    if (error !== 'cancel') { // 用户点击了取消按钮
      console.error('Failed to cancel reservation:', error)
      let message = 'Failed to cancel reservation'
      if (error.response?.data?.error?.message) {
        message = error.response.data.error.message
      } else if (error.response?.data?.message) {
        message = error.response.data.message
      }
      ElMessage.error(message)
    }
  } finally {
    isCancelling.value = false
  }
}

// 获取电影标题
const getMovieTitle = (showtime) => {
  if (!showtime) return 'N/A'
  // 如果showtime有movie对象，直接返回
  if (showtime.movie?.title) return showtime.movie.title
  // 否则返回电影ID
  return `Movie ID: ${showtime.movie_id || 'Unknown'}`
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 获取角色标签类型
const getRoleTagType = (role) => {
  switch (role) {
    case 'admin':
      return 'danger'
    case 'user':
      return 'success'
    default:
      return 'info'
  }
}

// 格式化时间
const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  if (!userStore.isAuthenticated) {
    ElMessage.error('Please login to view your profile')
    // 这里应该跳转到登录页，但为了演示，我们仅显示错误
    return
  }
  fetchReservations()
})
</script>

<style scoped>
.profile-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-avatar {
  width: 80px;
  height: 80px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 50%;
}

.user-info h2 {
  margin: 0 0 5px 0;
}

.user-info p {
  margin: 0;
  color: #666;
}

.reservations-card h3 {
  margin-top: 0;
}

.loading-container {
  margin: 20px 0;
}

.reservations-list {
  margin-top: 20px;
}
</style>
