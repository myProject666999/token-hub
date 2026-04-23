<template>
  <el-container class="user-layout">
    <el-header class="header">
      <div class="logo" @click="$router.push('/')">
        <el-icon><Connection /></el-icon>
        <span>Token Hub</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        mode="horizontal"
        background-color="transparent"
        text-color="#333"
        active-text-color="#409eff"
        class="nav-menu"
      >
        <el-menu-item index="/user/dashboard" @click="$router.push('/user/dashboard')">
          <el-icon><DataAnalysis /></el-icon>
          <span>控制台</span>
        </el-menu-item>
        <el-menu-item index="/user/recharge" @click="$router.push('/user/recharge')">
          <el-icon><Money /></el-icon>
          <span>充值中心</span>
        </el-menu-item>
        <el-menu-item index="/user/api-keys" @click="$router.push('/user/api-keys')">
          <el-icon><Key /></el-icon>
          <span>API密钥</span>
        </el-menu-item>
        <el-menu-item index="/user/call-logs" @click="$router.push('/user/call-logs')">
          <el-icon><Document /></el-icon>
          <span>调用日志</span>
        </el-menu-item>
      </el-menu>
      <div class="user-info">
        <el-dropdown @command="handleCommand">
          <span class="dropdown-trigger">
            <el-avatar :size="32" class="avatar">
              <el-icon><UserFilled /></el-icon>
            </el-avatar>
            <span class="username">{{ userStore.userInfo?.username }}</span>
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                个人信息
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleCommand = (command) => {
  if (command === 'profile') {
    router.push('/user/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/')
  }
}

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.fetchUserInfo()
  }
})
</script>

<style scoped lang="scss">
.user-layout {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.header {
  display: flex;
  align-items: center;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0 40px;

  .logo {
    display: flex;
    align-items: center;
    font-size: 20px;
    font-weight: bold;
    color: #409eff;
    cursor: pointer;

    .el-icon {
      font-size: 24px;
      margin-right: 8px;
    }
  }

  .nav-menu {
    flex: 1;
    margin-left: 40px;
    border-bottom: none;
  }

  .user-info {
    .dropdown-trigger {
      display: flex;
      align-items: center;
      cursor: pointer;

      .avatar {
        background-color: #409eff;
      }

      .username {
        margin: 0 8px;
      }
    }
  }
}

.main {
  padding: 20px 40px;
}
</style>
