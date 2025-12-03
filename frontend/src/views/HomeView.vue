<script setup>
import { ref, onMounted } from 'vue'
import { movieAPI } from '@/api/movie'
import { useRouter } from 'vue-router'
import { Search, Film } from '@element-plus/icons-vue'

const movies = ref([])
const filteredMovies = ref([])
const isLoading = ref(true)
const error = ref(null)
const searchQuery = ref('')
const router = useRouter()

onMounted(async () => {
  try {
    const response = await movieAPI.getMovies()
    movies.value = response.data.data || []
    filteredMovies.value = response.data.data || []
  } catch (err) {
    console.error('Failed to fetch movies:', err)
    // 检查是否是认证错误
    if (err.response?.status === 401) {
      error.value = 'Please login to view movies.'
    } else {
      error.value = 'Failed to load movies. Please try again later.'
    }
  } finally {
    isLoading.value = false
  }
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
</script>

<template>
  <div>
    <div class="search-container">
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
          Please <router-link to="/login" style="color: #409EFF; text-decoration: underline;">login</router-link> to view movies.
        </div>
      </template>
    </el-alert>

    <div v-if="!isLoading && !error">
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
  </div>
</template>

<style scoped>
.search-container {
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