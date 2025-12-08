<script setup>
import { RouterView, RouterLink, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { onMounted } from 'vue'
import { User, SwitchButton, House, Plus } from '@element-plus/icons-vue'

const userStore = useUserStore()
const router = useRouter()

onMounted(() => {
  userStore.checkAuth()
})

const handleLogout = async () => {
  await userStore.logout()
  router.push('/login')
}
</script>

<template>
  <el-container class="main-container">
    <el-header class="header">
      <div class="header-content">
        <RouterLink to="/" class="logo">
          <el-icon :size="28" style="margin-right: 8px;">
            <House />
          </el-icon>
          MovieHub
        </RouterLink>
        <el-menu mode="horizontal" :router="true" class="nav-menu">
          <el-menu-item index="/">Home</el-menu-item>

          <template v-if="!userStore.isAuthenticated">
            <el-menu-item index="/login">
              <el-icon><User /></el-icon>
              Login
            </el-menu-item>
            <el-menu-item index="/register">
              <el-icon><SwitchButton /></el-icon>
              Register
            </el-menu-item>
          </template>

          <template v-else>
            <!-- 管理员功能 -->
            <el-menu-item v-if="userStore.user?.role === 'admin'" index="/admin/add-movie">
              <el-icon><Plus /></el-icon>
              Add Movie
            </el-menu-item>
            <el-menu-item v-if="userStore.user?.role === 'admin'" index="/admin/halls">
              <el-icon><House /></el-icon>
              Manage Halls
            </el-menu-item>

            <el-menu-item index="/profile">
              <el-icon><User /></el-icon>
              Profile
            </el-menu-item>
            <el-menu-item @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>
              Logout
            </el-menu-item>
          </template>
        </el-menu>
      </div>
    </el-header>
    <el-main class="main-content">
      <RouterView />
    </el-main>
    <el-footer class="footer">
      <p>© 2025 MovieHub. All rights reserved.</p>
    </el-footer>
  </el-container>
</template>

<style scoped>
.main-container {
  min-height: 100vh;
}

.header {
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.logo {
  font-size: 24px;
  font-weight: bold;
  text-decoration: none;
  color: #409EFF;
  display: flex;
  align-items: center;
}

.nav-menu {
  border-bottom: none;
  background-color: transparent;
}

.main-content {
  padding: 20px;
  background-color: #f5f7fa;
}

.footer {
  background-color: #f5f5f5;
  text-align: center;
  padding: 20px;
  color: #666;
}
</style>
