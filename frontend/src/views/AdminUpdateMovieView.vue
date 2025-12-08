<template>
  <div class="update-movie-container">
    <el-card class="update-movie-card">
      <h2>Update Movie</h2>
      
      <el-form
        :model="movieForm"
        :rules="movieRules"
        ref="updateMovieFormRef"
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
            @click="handleUpdateMovie"
            :loading="isLoading"
            style="width: 100%"
          >
            Update Movie
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            type="default"
            size="large"
            @click="handleBack"
            style="width: 100%"
          >
            Back to Movie Detail
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { movieAPI } from '@/api/movie'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const isLoading = ref(false)
const updateMovieFormRef = ref()

const movieForm = reactive({
  title: '',
  description: ''
})

const movieRules = {
  title: [
    { required: true, message: 'Please enter a movie title', trigger: 'blur' },
    { min: 1, max: 100, message: 'Title must be between 1 and 100 characters', trigger: 'blur' }
  ]
}

// 获取电影详情
const fetchMovieDetail = async () => {
  try {
    const response = await movieAPI.getMovieById(route.params.id)
    const movie = response.data.data
    movieForm.title = movie.title
    movieForm.description = movie.description
  } catch (error) {
    console.error('Failed to fetch movie:', error)
    let message = 'Failed to load movie details'
    if (error.response?.data?.error?.message) {
      message = error.response.data.error.message
    } else if (error.response?.data?.message) {
      message = error.response.data.message
    }
    ElMessage.error(message)
    router.push('/').catch(() => {})
  }
}

const handleUpdateMovie = async () => {
  try {
    await updateMovieFormRef.value.validate()
    isLoading.value = true

    await movieAPI.updateMovie(route.params.id, movieForm.title, movieForm.description)
    ElMessage.success('Movie updated successfully!')

    // Redirect back to movie detail after successful update
    router.push(`/movies/${route.params.id}`)
  } catch (error) {
    console.error('Update movie error:', error)
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('Failed to update movie. Please try again.')
    }
  } finally {
    isLoading.value = false
  }
}

const handleBack = () => {
  router.push(`/movies/${route.params.id}`)
}

onMounted(() => {
  fetchMovieDetail()
})
</script>

<style scoped>
.update-movie-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.update-movie-card {
  padding: 30px;
}
</style>