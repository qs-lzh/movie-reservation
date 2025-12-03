<template>
  <div class="seat-reservation-container">
    <el-page-header @back="goBack" content="Reserve Showtime" />

    <el-card class="info-card" v-if="showtime">
      <div class="info-content">
        <div class="movie-info">
          <h2>{{ movieTitle }}</h2>
          <p>{{ formatDate(showtime.start_at) }} at {{ formatTime(showtime.start_at) }}</p>
          <p>Hall {{ showtime.hall_id }}</p>
        </div>
        <div class="availability-info">
          <el-tag :type="availabilityTagType" size="large">
            {{ remainingTickets }} tickets available
          </el-tag>
        </div>
      </div>
    </el-card>

    <el-card class="reservation-card">
      <h3>Reservation Information</h3>
      <el-alert
        title="Note"
        type="info"
        description="Each reservation is for one ticket. You can make multiple reservations if needed."
        :closable="false"
        style="margin-bottom: 20px;"
      />

      <div class="reservation-summary">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="Movie">{{ movieTitle }}</el-descriptions-item>
          <el-descriptions-item label="Date">{{ showtime ? formatDate(showtime.start_at) : 'Loading...' }}</el-descriptions-item>
          <el-descriptions-item label="Time">{{ showtime ? formatTime(showtime.start_at) : 'Loading...' }}</el-descriptions-item>
          <el-descriptions-item label="Hall">Hall {{ showtime?.hall_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Available Tickets">{{ remainingTickets }}</el-descriptions-item>
        </el-descriptions>

        <el-button
          type="primary"
          size="large"
          :disabled="remainingTickets === 0 || isReserving"
          @click="confirmReservation"
          :loading="isReserving"
          style="width: 100%; margin-top: 30px;"
        >
          {{ remainingTickets === 0 ? 'Sold Out' : 'Confirm Reservation (1 ticket)' }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showtimeAPI } from '@/api/showtime'
import { movieAPI } from '@/api/movie'
import { reservationAPI } from '@/api/reservation'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const showtime = ref(null)
const movieTitle = ref('Loading...')
const remainingTickets = ref(0)
const isReserving = ref(false)

// 获取场次信息
const fetchShowtimeInfo = async () => {
  try {
    // 获取场次详情
    const showtimeResponse = await showtimeAPI.getShowtimeById(route.params.id)
    showtime.value = showtimeResponse.data.data

    // 获取电影信息
    if (showtime.value?.movie_id) {
      const movieResponse = await movieAPI.getMovieById(showtime.value.movie_id)
      movieTitle.value = movieResponse.data.data?.title || 'Unknown Movie'
    }

    // 获取场次的可用性
    const availabilityResponse = await showtimeAPI.getShowtimeAvailability(route.params.id)
    remainingTickets.value = availabilityResponse.data.data.remaining_tickets
  } catch (error) {
    console.error('Failed to fetch showtime info:', error)
    ElMessage.error('Failed to load showtime information')
    router.push('/')
  }
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 格式化时间
const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 确认预订
const confirmReservation = async () => {
  try {
    // 确认对话框
    await ElMessageBox.confirm(
      'You are about to reserve 1 ticket for this showtime. Continue?',
      'Confirm Reservation',
      {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'info'
      }
    )

    isReserving.value = true

    // 调用预订API - 只需要场次ID
    await reservationAPI.createReservation(route.params.id)

    ElMessage.success('Reservation confirmed successfully!')
    router.push('/profile') // 跳转到用户个人页面查看预订
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了预订
      return
    }
    console.error('Reservation failed:', error)
    let message = 'Reservation failed'
    if (error.response?.data?.message) {
      message = error.response.data.message
    } else if (error.response?.data?.code) {
      message = error.response.data.code
    }
    ElMessage.error(message)
  } finally {
    isReserving.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 可用性标签类型
const availabilityTagType = computed(() => {
  if (remainingTickets.value === 0) return 'danger'
  if (remainingTickets.value < 10) return 'warning'
  return 'success'
})

onMounted(() => {
  fetchShowtimeInfo()
})
</script>

<style scoped>
.seat-reservation-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.info-card {
  margin-bottom: 20px;
}

.info-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.movie-info h2 {
  margin: 0 0 10px 0;
}

.movie-info p {
  margin: 5px 0;
  color: #606266;
}

.availability-info {
  text-align: right;
}

.reservation-card {
  padding: 30px;
}

.reservation-card h3 {
  margin-top: 0;
  margin-bottom: 20px;
}

.reservation-summary {
  margin-top: 20px;
}
</style>
