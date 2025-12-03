<template>
  <div class="add-movie-container">
    <el-card class="add-movie-card">
      <h2>Admin Panel</h2>
      <el-tabs v-model="activeTab" class="admin-tabs">
        <!-- 添加电影 -->
        <el-tab-pane label="Add Movie" name="movie">
          <el-form
            :model="movieForm"
            :rules="movieRules"
            ref="addMovieFormRef"
            label-width="120px"
          >
            <el-form-item label="Title" prop="title">
              <el-input
                v-model="movieForm.title"
                placeholder="Enter movie title"
                size="large"
              />
            </el-form-item>
            <el-form-item label="Description" prop="description">
              <el-input
                v-model="movieForm.description"
                type="textarea"
                :rows="4"
                placeholder="Enter movie description"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                @click="handleAddMovie"
                :loading="isLoading"
                style="width: 100%"
              >
                Add Movie
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 添加场次 -->
        <el-tab-pane label="Add Showtime" name="showtime">
          <el-form
            :model="showtimeForm"
            :rules="showtimeRules"
            ref="addShowtimeFormRef"
            label-width="120px"
          >
            <el-form-item label="Movie" prop="movieId">
              <el-select
                v-model="showtimeForm.movieId"
                placeholder="Select a movie"
                size="large"
                style="width: 100%"
                filterable
              >
                <el-option
                  v-for="movie in movies"
                  :key="movie.id"
                  :label="movie.title"
                  :value="movie.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="Date & Time" prop="startAt">
              <el-date-picker
                v-model="showtimeForm.startAt"
                type="datetime"
                placeholder="Select date and time"
                size="large"
                style="width: 100%"
                :disabled-date="disabledDate"
              />
            </el-form-item>
            <el-form-item label="Hall ID" prop="hallId">
              <el-input-number
                v-model="showtimeForm.hallId"
                :min="1"
                :max="10"
                size="large"
                style="width: 100%"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                @click="handleAddShowtime"
                :loading="isLoading"
                style="width: 100%"
              >
                Add Showtime
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { movieAPI } from '@/api/movie'
import { showtimeAPI } from '@/api/showtime'
import { ElMessage } from 'element-plus'

const router = useRouter()
const isLoading = ref(false)
const addMovieFormRef = ref()
const addShowtimeFormRef = ref()
const activeTab = ref('movie')
const movies = ref([])

const movieForm = reactive({
  title: '',
  description: ''
})

const showtimeForm = reactive({
  movieId: null,
  startAt: null,
  hallId: 1
})

const movieRules = {
  title: [
    { required: true, message: 'Please enter a movie title', trigger: 'blur' },
    { min: 1, max: 100, message: 'Title must be between 1 and 100 characters', trigger: 'blur' }
  ]
}

const showtimeRules = {
  movieId: [
    { required: true, message: 'Please select a movie', trigger: 'change' }
  ],
  startAt: [
    { required: true, message: 'Please select date and time', trigger: 'change' }
  ],
  hallId: [
    { required: true, message: 'Please enter hall ID', trigger: 'blur' }
  ]
}

// 禁用过去的日期
const disabledDate = (time) => {
  return time.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

// 获取所有电影列表
const fetchMovies = async () => {
  try {
    const response = await movieAPI.getMovies()
    movies.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch movies:', error)
    let message = 'Failed to load movies'
    if (error.response?.data?.error?.message) {
      message = error.response.data.error.message
    } else if (error.response?.data?.message) {
      message = error.response.data.message
    }
    ElMessage.error(message)
  }
}

const handleAddMovie = async () => {
  try {
    await addMovieFormRef.value.validate()
    isLoading.value = true

    await movieAPI.createMovie(movieForm.title, movieForm.description)
    ElMessage.success('Movie added successfully!')

    // 重置表单
    movieForm.title = ''
    movieForm.description = ''

    // 刷新电影列表
    await fetchMovies()

    // 切换到添加场次标签
    activeTab.value = 'showtime'
  } catch (error) {
    console.error('Add movie error:', error)
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('Failed to add movie. Please try again.')
    }
  } finally {
    isLoading.value = false
  }
}

const handleAddShowtime = async () => {
  try {
    await addShowtimeFormRef.value.validate()
    isLoading.value = true

    await showtimeAPI.createShowtime(
      showtimeForm.movieId,
      showtimeForm.startAt.toISOString(),
      showtimeForm.hallId
    )
    ElMessage.success('Showtime added successfully!')

    // 重置表单
    showtimeForm.movieId = null
    showtimeForm.startAt = null
    showtimeForm.hallId = 1
  } catch (error) {
    console.error('Add showtime error:', error)
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('Failed to add showtime. Please try again.')
    }
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchMovies()
})
</script>

<style scoped>
.add-movie-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.add-movie-card {
  padding: 30px;
}

.admin-tabs {
  margin-top: 20px;
}
</style>