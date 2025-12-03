<template>
  <div class="movie-detail-container">
    <el-card v-if="movie" class="movie-card">
      <div class="movie-header">
        <div class="movie-poster-placeholder">
          <el-icon :size="80" class="movie-poster-icon">
            <Film />
          </el-icon>
        </div>
        <div class="movie-info">
          <h1 class="movie-title">{{ movie.title }}</h1>
          <p class="movie-description" v-if="movie.description">{{ movie.description }}</p>
        </div>
      </div>
    </el-card>

    <el-card class="showtimes-container" v-if="showtimes.length > 0">
      <h3>Showtimes</h3>
      <el-row :gutter="20">
        <el-col
          v-for="showtime in showtimes"
          :key="showtime.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card
            class="showtime-card"
            :body-style="{ padding: '20px' }"
            @click="goToReservation(showtime.id)"
          >
            <div class="showtime-info">
              <div class="showtime-date">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(showtime.start_at) }}
              </div>
              <div class="showtime-time">
                <el-icon><Clock /></el-icon>
                {{ formatTime(showtime.start_at) }}
              </div>
              <div class="showtime-hall">
                <el-icon><Location /></el-icon>
                Hall {{ showtime.hall_id }}
              </div>
              <div class="showtime-availability">
                <el-icon><Tickets /></el-icon>
                <span v-if="showtime.availability !== undefined">
                  {{ showtime.availability }} seats available
                </span>
                <span v-else>Loading...</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <el-empty
      v-else-if="!isLoading && showtimes.length === 0"
      description="No showtimes available"
    />

    <div v-if="isLoading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { movieAPI } from '@/api/movie'
import { showtimeAPI } from '@/api/showtime'
import { ElMessage } from 'element-plus'
import { Film, Calendar, Clock, Location, Tickets } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const movie = ref(null)
const showtimes = ref([])
const isLoading = ref(true)

// 获取电影详情
const fetchMovieDetail = async () => {
  try {
    const response = await movieAPI.getMovieById(route.params.id)
    movie.value = response.data.data
    await fetchShowtimes()
  } catch (error) {
    console.error('Failed to fetch movie:', error)
    ElMessage.error('Failed to load movie details')
  }
}

// 获取场次信息
const fetchShowtimes = async () => {
  try {
    const response = await movieAPI.getMovieShowtimes(route.params.id)
    showtimes.value = response.data.data || []

    // 获取每个场次的可用性信息
    for (const showtime of showtimes.value) {
      try {
        const availabilityResponse = await showtimeAPI.getShowtimeAvailability(showtime.id)
        showtime.availability = availabilityResponse.data.data.remaining_tickets
      } catch (error) {
        console.error(`Failed to fetch availability for showtime ${showtime.id}:`, error)
        showtime.availability = 'Error'
      }
    }
  } catch (error) {
    console.error('Failed to fetch showtimes:', error)
    ElMessage.error('Failed to load showtimes')
  } finally {
    isLoading.value = false
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

// 跳转到预订页面
const goToReservation = (showtimeId) => {
  router.push({ name: 'seat-reservation', params: { id: showtimeId } })
}

onMounted(() => {
  fetchMovieDetail()
})
</script>

<style scoped>
.movie-detail-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.movie-card {
  margin-bottom: 30px;
}

.movie-header {
  display: flex;
  gap: 30px;
}

.movie-poster-placeholder {
  width: 200px;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 8px;
  flex-shrink: 0;
}

.movie-poster-icon {
  color: #c0c4cc;
}

.movie-info {
  flex: 1;
}

.movie-title {
  margin: 0 0 15px 0;
  font-size: 2em;
  color: #303133;
}

.movie-description {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  font-size: 1.1em;
}

.showtimes-container {
  margin-top: 20px;
}

.showtimes-container h3 {
  margin: 0 0 20px 0;
  font-size: 1.5em;
}

.showtime-card {
  cursor: pointer;
  transition: transform 0.3s ease;
  margin-bottom: 20px;
}

.showtime-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.showtime-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.showtime-date, .showtime-time, .showtime-hall, .showtime-availability {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9em;
}

.showtime-date, .showtime-time {
  font-weight: 500;
}

.loading-container {
  margin-top: 30px;
}
</style>
