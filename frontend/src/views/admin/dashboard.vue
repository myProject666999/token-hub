<template>
  <div class="admin-dashboard">
    <h2 class="page-title">数据概览</h2>

    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card users-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.total_users || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card calls-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.total_calls || 0 }}</div>
              <div class="stat-label">总调用次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card recharge-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ overview?.total_recharge || 0 }}</div>
              <div class="stat-label">总充值金额</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card points-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Coin /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.total_points_used || 0 }}</div>
              <div class="stat-label">总消耗积分</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>今日数据</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="今日调用次数">
              <span class="highlight">{{ overview?.today_calls || 0 }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="今日充值金额">
              <span class="highlight highlight-green">¥{{ overview?.today_recharge || 0 }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/admin/providers')">
              <el-icon><Cpu /></el-icon>
              管理服务商
            </el-button>
            <el-button type="success" @click="$router.push('/admin/models')">
              <el-icon><Connection /></el-icon>
              管理模型
            </el-button>
            <el-button type="warning" @click="$router.push('/admin/users')">
              <el-icon><User /></el-icon>
              用户管理
            </el-button>
            <el-button type="info" @click="$router.push('/admin/points-config')">
              <el-icon><Setting /></el-icon>
              积分配置
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOverviewStatistics } from '@/api/user'

const overview = ref(null)

const loadData = async () => {
  try {
    const res = await getOverviewStatistics()
    overview.value = res.data
  } catch (e) {
    console.error('Failed to load overview:', e)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.admin-dashboard {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }

  .stats-row {
    margin-bottom: 24px;
  }

  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
    }

    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 28px;
      color: #fff;
      margin-right: 16px;
    }

    .stat-info {
      flex: 1;
    }

    .stat-value {
      font-size: 28px;
      font-weight: bold;
      color: #303133;
    }

    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 4px;
    }
  }

  .users-card .stat-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .calls-card .stat-icon {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  }

  .recharge-card .stat-icon {
    background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
  }

  .points-card .stat-icon {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  }

  .section-card {
    margin-bottom: 24px;
  }

  .card-header {
    font-weight: bold;
    font-size: 16px;
  }

  .highlight {
    color: #409eff;
    font-weight: bold;
    font-size: 18px;
  }

  .highlight-green {
    color: #67c23a;
  }

  .quick-actions {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;

    .el-button {
      width: 100%;
      height: 50px;
    }
  }
}
</style>
