<template>
  <div class="register-container">
    <el-card class="register-form">
      <h2>Create an Account</h2>
      <el-form :model="form" :rules="rules" ref="registerFormRef">
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
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="Confirm Password"
            size="large"
            prefix-icon="Lock"
          />
        </el-form-item>
        <el-form-item label="Role">
          <el-radio-group v-model="form.role">
            <el-radio label="user">User</el-radio>
            <el-radio label="admin">Admin</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item prop="adminRolePassword" v-if="form.role === 'admin'">
          <el-input
            v-model="form.adminRolePassword"
            type="password"
            placeholder="Admin Role Password"
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
            Register
          </el-button>
        </el-form-item>
      </el-form>
      <p class="login-link">
        Already have an account?
        <RouterLink to="/login">Login here</RouterLink>
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
const registerFormRef = ref()
const captchaRef = ref(null);

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  role: 'user',
  adminRolePassword: "",
})

const rules = {
  username: [
    { required: true, message: 'Please enter a username', trigger: 'blur' },
    { min: 3, message: 'Username must be at least 3 characters', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter a password', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: 'Please confirm your password', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== form.password) {
          callback(new Error('Passwords do not match'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const handleSubmit = async () => {
  try {
    await registerFormRef.value.validate()

    // Make sure captcha is completed before registration
    if (!key.value) {
      ElMessage.error('Please complete the captcha verification')
      return
    }

    isLoading.value = true

    const result = await userStore.register(form.username, form.password, form.role, key.value,
    form.adminRolePassword)

    if (result.success) {
      ElMessage.success(result.message)
      router.push('/login')
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    console.error('Registration error:', error)
    ElMessage.error('Registration failed. Please try again.')
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
    clickedDots = [];
  }
};

function refreshCaptcha() {
  captchaRef.value.refresh();
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
}

.register-form {
  width: 100%;
  max-width: 400px;
  padding: 30px;
}

.login-link {
  text-align: center;
  margin-top: 20px;
}

.login-link a {
  color: #409EFF;
  text-decoration: none;
}
</style>
