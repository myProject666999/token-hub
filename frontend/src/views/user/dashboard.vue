<template>
  <div class="dashboard-page">
    <h2 class="page-title">控制台</h2>

    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card points-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Coin /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ userPoints?.points || 0 }}</div>
              <div class="stat-label">可用积分</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card total-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ userPoints?.total_points || 0 }}</div>
              <div class="stat-label">累计获得</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card used-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ userPoints?.used_points || 0 }}</div>
              <div class="stat-label">已使用</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card call-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ callStats?.total_calls || 0 }}</div>
              <div class="stat-label">总调用次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="section-card">
          <template #header>
            <div class="card-header">
              <span>使用统计</span>
            </div>
          </template>
          <div class="stats-detail">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="总调用次数">
                {{ callStats?.total_calls || 0 }}
              </el-descriptions-item>
              <el-descriptions-item label="成功次数">
                <el-tag type="success">{{ callStats?.success_calls || 0 }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="失败次数">
                <el-tag type="danger">{{ callStats?.failed_calls || 0 }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="总消耗积分">
                <span class="highlight">{{ callStats?.total_points || 0 }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="输入Token总数">
                {{ callStats?.total_input_tokens || 0 }}
              </el-descriptions-item>
              <el-descriptions-item label="输出Token总数">
                {{ callStats?.total_output_tokens || 0 }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="section-card quick-actions">
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>
          <el-button type="primary" class="action-btn" @click="$router.push('/user/recharge')">
            <el-icon><Money /></el-icon>
            立即充值
          </el-button>
          <el-button class="action-btn" @click="$router.push('/user/api-keys')">
            <el-icon><Key /></el-icon>
            管理API密钥
          </el-button>
          <el-button class="action-btn" @click="$router.push('/user/call-logs')">
            <el-icon><Document /></el-icon>
            查看调用日志
          </el-button>
          <el-button class="action-btn" @click="$router.push('/docs')">
            <el-icon><Reading /></el-icon>
            查看接入文档
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserPoints, getUserStatistics } from '@/api/auth'
import { getUserCallStatistics } from '@/api/user'

const userPoints = ref(null)
const callStats = ref(null)

const loadData = async () => {
  try {
    const [pointsRes, statsRes] = await Promise.all([
      getUserPoints(),
      getUserCallStatistics()
    ])
    userPoints.value = pointsRes.data
    callStats.value = statsRes.data
  } catch (e) {
    console.error('Failed to load dashboard data:', e)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.dashboard-page {
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

  .points-card .stat-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .total-card .stat-icon {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  }

  .used-card .stat-icon {
    background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
  }

  .call-card .stat-icon {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  }

  .section-card {
    margin-bottom: 24px;
  }

  .card-header {
    font-weight: bold;
    font-size: 16px;
  }

  .stats-detail {
    .highlight {
      color: #409eff;
      font-weight: bold;
    }
  }

  .quick-actions {
    .action-btn {
      width: 100%;
      margin-bottom: 12px;
      height: 48px;
    }
  }
}
</style>
