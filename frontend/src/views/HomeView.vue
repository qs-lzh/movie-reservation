<script setup>
import { ref, onMounted } from 'vue'
import { movieAPI } from '@/api/movie'
import { showtimeAPI } from '@/api/showtime'
import { useRouter } from 'vue-router'
import { Search, Film, Calendar } from '@element-plus/icons-vue'

const movies = ref([])
const filteredMovies = ref([])
const showtimes = ref([])
const isLoading = ref(true)
const error = ref(null)
const searchQuery = ref('')
const currentView = ref('movies') // 'movies' or 'showtimes'
const router = useRouter()

const fetchAllData = async () => {
  let moviesError = null;
  let showtimesError = null;

  try {
    // 并行获取电影和场次数据
    const [moviesResponse, showtimesResponse] = await Promise.allSettled([
      movieAPI.getMovies(),
      showtimeAPI.getShowtimes()
    ])

    if (moviesResponse.status === 'fulfilled') {
      movies.value = moviesResponse.value.data.data || []
      filteredMovies.value = moviesResponse.value.data.data || []
    } else {
      console.error('Failed to fetch movies:', moviesResponse.reason)
      moviesError = moviesResponse.reason;
    }

    if (showtimesResponse.status === 'fulfilled') {
      showtimes.value = showtimesResponse.value.data.data || []
      // 为每个场次添加movie_title属性
      for (const showtime of showtimes.value) {
        if (!showtime.movie_title && showtime.movie_id) {
          const movie = movies.value.find(m => m.id === showtime.movie_id)
          if (movie) {
            showtime.movie_title = movie.title
          }
        }
      }
      // 按时间排序场次，最早的在前
      showtimes.value.sort((a, b) => new Date(a.start_at) - new Date(b.start_at))
    } else {
      console.error('Failed to fetch showtimes:', showtimesResponse.reason)
      showtimesError = showtimesResponse.reason;
    }

    // 确定总体错误状态
    if (moviesError && showtimesError) {
      // 两个请求都失败，检查是否有认证错误
      if (moviesError?.response?.status === 401 || showtimesError?.response?.status === 401) {
        error.value = 'Please login to view content.'
      } else {
        error.value = 'Failed to load movies and showtimes. Please try again later.'
      }
    } else if (moviesError) {
      // 仅电影请求失败
      if (moviesError?.response?.status === 401) {
        error.value = 'Please login to view content.'
      } else {
        error.value = 'Failed to load movies. Please try again later.'
      }
    } else if (showtimesError) {
      // 仅场次请求失败
      if (showtimesError?.response?.status === 401) {
        error.value = 'Please login to view content.'
      } else {
        error.value = 'Failed to load showtimes. Please try again later.'
      }
    }
    // 如果都没有错误，继续正常显示

  } catch (err) {
    console.error('Error fetching data:', err)
    // 检查是否是认证错误
    if (err.response?.status === 401) {
      error.value = 'Please login to view content.'
    } else {
      error.value = 'Failed to load content. Please try again later.'
    }
  } finally {
    isLoading.value = false
  }
}

onMounted(async () => {
  await fetchAllData()
})

// 搜索过滤
const filterMovies = () => {
  if (!searchQuery.value) {
    filteredMovies.value = movies.value
  } else {
    filteredMovies.value = movies.value.filter(movie =>
      movie.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      movie.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
}

const goToMovieDetail = (id) => {
  router.push({ name: 'movie-detail', params: { id } })
}

const goToReservation = (showtimeId) => {
  router.push({ name: 'seat-reservation', params: { id: showtimeId } })
}

// 切换视图
const switchView = (view) => {
  currentView.value = view
  if (view === 'movies') {
    searchQuery.value = '' // 清空搜索查询，因为我们只对电影进行搜索
  }
}

// 获取场次对应的电影标题
const getShowtimeMovieTitle = (showtime) => {
  // 首先尝试从已获取的电影列表中找到标题
  if (showtime.movie_title) {
    return showtime.movie_title
  } else if (showtime.movie_id) {
    // 如果没有movie_title但有movie_id，尝试从本地电影列表查找
    const movie = movies.value.find(m => m.id === showtime.movie_id)
    if (movie && movie.title) {
      return movie.title
    }
    // 如果仍然找不到，显示ID
    return `Movie ID: ${showtime.movie_id}`
  }
  return 'Unknown Movie'
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

// 格式化时间
const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取可用性标签类型
const getAvailabilityTagType = (availableTickets) => {
  if (availableTickets === 0) return 'danger'
  if (availableTickets < 10) return 'warning'
  return 'success'
}
</script>

<template>
  <div>
    <!-- 视图切换按钮 -->
    <div class="view-toggle-container">
      <el-button-group>
        <el-button
          :type="currentView === 'movies' ? 'primary' : 'default'"
          @click="switchView('movies')"
        >
          <el-icon><Film /></el-icon>
          Movies
        </el-button>
        <el-button
          :type="currentView === 'showtimes' ? 'primary' : 'default'"
          @click="switchView('showtimes')"
        >
          <el-icon><Calendar /></el-icon>
          Showtimes
        </el-button>
      </el-button-group>
    </div>

    <!-- 电影视图的搜索框 -->
    <div v-if="currentView === 'movies'" class="search-container">
      <el-input
        v-model="searchQuery"
        placeholder="Search movies..."
        size="large"
        @input="filterMovies"
        clearable
        style="max-width: 500px;"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <div v-if="isLoading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      :closable="false"
    >
      <template v-if="error.includes('login')" #default>
        <div>
          Please <router-link to="/login" style="color: #409EFF; text-decoration: underline;">login</router-link> to view content.
        </div>
      </template>
    </el-alert>

    <div v-if="!isLoading && !error">
      <!-- 电影视图 -->
      <div v-if="currentView === 'movies'">
        <el-row :gutter="20">
          <el-col
            v-for="movie in filteredMovies"
            :key="movie.id"
            :xs="24"
            :sm="12"
            :md="8"
            :lg="6"
          >
            <el-card
              shadow="hover"
              class="movie-card"
              @click="goToMovieDetail(movie.id)"
            >
              <div class="movie-poster-placeholder">
                <el-icon :size="60" class="movie-poster-icon">
                  <Film />
                </el-icon>
              </div>
              <div class="movie-info">
                <h3 class="movie-title">{{ movie.title }}</h3>
                <p class="movie-description" v-if="movie.description">
                  {{ movie.description.substring(0, 100) }}{{ movie.description.length > 100 ? '...' : '' }}
                </p>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-empty
          v-if="filteredMovies.length === 0"
          description="No movies found"
        />
      </div>

      <!-- 场次视图 -->
      <div v-else-if="currentView === 'showtimes'">
        <el-table
          :data="showtimes"
          style="width: 100%"
          stripe
        >
          <el-table-column prop="id" label="Showtime ID" width="120" />
          <el-table-column label="Movie" width="200">
            <template #default="{ row }">
              {{ getShowtimeMovieTitle(row) }}
            </template>
          </el-table-column>
          <el-table-column label="Date & Time" width="200">
            <template #default="{ row }">
              {{ formatDate(row.start_at) }}<br>
              {{ formatTime(row.start_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="hall_id" label="Hall" width="100" />
          <el-table-column label="Actions" width="150">
            <template #default="{ row }">
              <el-button
                size="small"
                type="primary"
                @click="goToReservation(row.id)"
              >
                Reserve
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-empty
          v-if="showtimes.length === 0"
          description="No showtimes available"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
  padding: 0 20px;
}

.view-toggle-container {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
  padding: 0 20px;
}

.movie-card {
  cursor: pointer;
  margin-bottom: 20px;
  transition: transform 0.3s ease;
}

.movie-card:hover {
  transform: translateY(-5px);
}

.movie-poster-placeholder {
  width: 100%;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.movie-poster-icon {
  color: #c0c4cc;
}

.movie-info {
  padding: 14px;
}

.movie-title {
  margin: 0 0 10px 0;
  font-size: 1.2em;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.movie-description {
  margin: 0;
  color: #666;
  font-size: 0.9em;
  line-height: 1.4;
}
</style>