<template>
  <div class="profile-page">
    <h2 class="page-title">个人信息</h2>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="avatar-card">
          <div class="avatar-section">
            <el-avatar :size="120" class="user-avatar">
              <el-icon size="60"><UserFilled /></el-icon>
            </el-avatar>
            <h3 class="username">{{ userInfo?.username }}</h3>
            <el-tag :type="userInfo?.role === 'admin' ? 'danger' : 'primary'">
              {{ userInfo?.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="info-card">
          <template #header>
            <span>基本信息</span>
          </template>
          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="100px"
          >
            <el-form-item label="用户名">
              <el-input v-model="profileForm.username" disabled />
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileForm.email" placeholder="请输入邮箱地址" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="注册时间">
              <el-input :value="userInfo?.created_at" disabled />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleUpdate">
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { getProfile, updateProfile } from '@/api/auth'

const userStore = useUserStore()
const profileFormRef = ref(null)
const loading = ref(false)

const userInfo = computed(() => userStore.userInfo)

const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: ''
})

const profileRules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const loadProfile = async () => {
  try {
    const res = await getProfile()
    profileForm.username = res.data.username
    profileForm.nickname = res.data.nickname || ''
    profileForm.email = res.data.email || ''
    profileForm.phone = res.data.phone || ''
  } catch (e) {
    console.error('Failed to load profile:', e)
  }
}

const handleUpdate = async () => {
  const valid = await profileFormRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone
    })
    ElMessage.success('更新成功')
    userStore.fetchUserInfo()
  } catch (e) {
    console.error('Failed to update profile:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (userStore.userInfo) {
    profileForm.username = userStore.userInfo.username
    profileForm.nickname = userStore.userInfo.nickname || ''
    profileForm.email = userStore.userInfo.email || ''
    profileForm.phone = userStore.userInfo.phone || ''
  }
  loadProfile()
})
</script>

<style scoped lang="scss">
.profile-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }

  .avatar-card {
    .avatar-section {
      text-align: center;
      padding: 20px 0;
    }

    .user-avatar {
      background-color: #409eff;
      margin-bottom: 16px;
    }

    .username {
      font-size: 20px;
      font-weight: bold;
      margin-bottom: 12px;
    }
  }

  .info-card {
    :deep(.el-form-item__label) {
      font-weight: 500;
    }
  }
}
</style>
