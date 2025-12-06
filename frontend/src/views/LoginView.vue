<template>
  <div class="login-container">
    <el-card class="login-form">
      <h2>Login to MovieHub</h2>
      <el-form :model="form" :rules="rules" ref="loginFormRef">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="Username"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="Password"
            size="large"
            prefix-icon="Lock"
          />
        </el-form-item>
        <gocaptcha-click
          :config="config"
          :data="data"
          :events="clickEvents"
          ref="captchaRef"
        />

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            @click="handleSubmit"
            :loading="isLoading"
            style="width: 100%"
          >
            Login
          </el-button>
        </el-form-item>
      </el-form>
      <p class="register-link">
        Don't have an account?
        <RouterLink to="/register">Register here</RouterLink>
      </p>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import apiClient from '../api/index.js';

const router = useRouter()
const userStore = useUserStore()
const isLoading = ref(false)
const loginFormRef = ref()
const captchaRef = ref(null);

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: 'Please enter your username', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter your password', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  try {
    await loginFormRef.value.validate()
    isLoading.value = true

    const result = await userStore.login(form.username, form.password)

    if (result.success) {
      ElMessage.success(result.message)
      router.push('/')
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    ElMessage.error('Login failed. Please check your credentials.')
  } finally {
    isLoading.value = false
  }
}

// 配置项
const config = {
  width: 300,
  height: 150,
  thumbWidth: 50,
  thumbHeight: 50,
  showTheme: true,
  title: "请完成验证"
};
const data = ref(null);
const key = ref(null);
async function fetchCaptcha() {
  try {
    const res = await apiClient.get("/captcha");
    data.value = {
      image: res.data.data.image,
      thumb: res.data.data.thumb,
    }
    key.value = res.data.data.key;
  } catch (err) {
    console.error("获取验证码失败:", err);
  }
}
onMounted(async () => {
  await fetchCaptcha();
  if (captchaRef.value) {
    captchaRef.value.refresh();
  }
});

let clickedDots = [];
const clickEvents = {
  click: (x, y) => {
    clickedDots.push({ x, y });
  },
  confirm: async (reset) => {
    try {
      const res = await apiClient.post("/captcha", { dots: clickedDots, key: key.value });

      if (res.data.data.success) {
        alert("验证成功");
        clickedDots = [];
        return true;
      } else {
        alert("验证失败");
        reset();
        clickedDots = [];
        return false;
      }
    } catch (err) {
      alert("提交失败，请重试");
      reset?.();
      clickedDots = [];
      return false;
    }
  },
  refresh: () => {
    clickedDots = []; // 刷新时清空点击点
  }
};

// 调用方法
function refreshCaptcha() {
  captchaRef.value.refresh();
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
}

.login-form {
  width: 100%;
  max-width: 400px;
  padding: 30px;
}

.register-link {
  text-align: center;
  margin-top: 20px;
}

.register-link a {
  color: #409EFF;
  text-decoration: none;
}
</style>
