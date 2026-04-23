<template>
  <div class="models-page">
    <header class="page-header">
      <div class="logo" @click="$router.push('/')">
        <el-icon><Connection /></el-icon>
        <span>Token Hub</span>
      </div>
      <nav class="nav">
        <el-button @click="$router.push('/')">首页</el-button>
        <el-button type="primary">模型列表</el-button>
        <el-button @click="$router.push('/docs')">接入说明</el-button>
        <template v-if="userStore.isAuthenticated">
          <el-button @click="$router.push('/user/dashboard')">控制台</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')">登录</el-button>
        </template>
      </nav>
    </header>

    <main class="main-content">
      <div class="page-title">
        <h1>可用大模型列表</h1>
        <p>选择合适的大模型，开启您的AI之旅</p>
      </div>

      <div class="filters">
        <el-select v-model="selectedProvider" placeholder="选择服务商" clearable @change="filterModels">
          <el-option
            v-for="provider in providers"
            :key="provider.id"
            :label="provider.name"
            :value="provider.id"
          />
        </el-select>
      </div>

      <div class="provider-section" v-for="provider in filteredProviders" :key="provider.id">
        <div class="provider-header">
          <h2>{{ provider.name }}</h2>
          <el-tag type="primary">{{ provider.models?.length || 0 }} 个模型</el-tag>
        </div>
        <p class="provider-desc">{{ provider.description }}</p>
        
        <div class="models-grid">
          <div class="model-card" v-for="model in provider.models" :key="model.id">
            <div class="model-header">
              <h3>{{ model.name }}</h3>
              <el-tag :type="model.status === 1 ? 'success' : 'danger'">
                {{ model.status === 1 ? '可用' : '维护中' }}
              </el-tag>
            </div>
            <p class="model-desc">{{ model.description }}</p>
            <div class="model-specs">
              <div class="spec-item">
                <span class="label">类型</span>
                <span class="value">{{ model.model_type || '文本' }}</span>
              </div>
              <div class="spec-item">
                <span class="label">最大Tokens</span>
                <span class="value">{{ model.max_tokens }}</span>
              </div>
              <div class="spec-item">
                <span class="label">上下文</span>
                <span class="value">{{ model.context_limit }}</span>
              </div>
            </div>
            <div class="model-pricing">
              <div class="pricing-item">
                <span class="label">输入价格</span>
                <span class="price">{{ model.points_per_1k_input }} 积分/千tokens</span>
              </div>
              <div class="pricing-item">
                <span class="label">输出价格</span>
                <span class="price">{{ model.points_per_1k_output }} 积分/千tokens</span>
              </div>
            </div>
            <div class="model-features">
              <el-tag v-if="model.supports_vision" size="small" type="info">支持图像</el-tag>
              <el-tag v-if="model.supports_function" size="small" type="warning">支持函数调用</el-tag>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { getProviders } from '@/api/provider'

const userStore = useUserStore()

const providers = ref([])
const selectedProvider = ref(null)

const filteredProviders = computed(() => {
  if (!selectedProvider.value) {
    return providers.value
  }
  return providers.value.filter(p => p.id === selectedProvider.value)
})

const filterModels = () => {
  // Filter is handled by computed
}

onMounted(async () => {
  try {
    const res = await getProviders(false)
    providers.value = res.data
  } catch (e) {
    console.error('Failed to load providers:', e)
  }
})
</script>

<style scoped lang="scss">
.models-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  height: 70px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;

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

  .nav {
    .el-button {
      margin-left: 10px;
    }
  }
}

.main-content {
  padding: 100px 40px 60px;
  max-width: 1400px;
  margin: 0 auto;

  .page-title {
    text-align: center;
    margin-bottom: 40px;

    h1 {
      font-size: 36px;
      margin-bottom: 10px;
    }

    p {
      color: #666;
      font-size: 16px;
    }
  }

  .filters {
    margin-bottom: 40px;

    .el-select {
      width: 200px;
    }
  }

  .provider-section {
    background: #fff;
    border-radius: 12px;
    padding: 30px;
    margin-bottom: 30px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

    .provider-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 10px;

      h2 {
        margin: 0;
        font-size: 24px;
      }
    }

    .provider-desc {
      color: #666;
      margin-bottom: 20px;
    }
  }

  .models-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;

    .model-card {
      border: 1px solid #e4e7ed;
      border-radius: 8px;
      padding: 20px;
      transition: box-shadow 0.3s;

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }

      .model-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 12px;

        h3 {
          margin: 0;
          font-size: 18px;
        }
      }

      .model-desc {
        color: #666;
        font-size: 14px;
        margin-bottom: 16px;
        line-height: 1.5;
      }

      .model-specs {
        display: flex;
        gap: 20px;
        margin-bottom: 16px;
        padding-bottom: 16px;
        border-bottom: 1px solid #e4e7ed;

        .spec-item {
          .label {
            display: block;
            font-size: 12px;
            color: #909399;
            margin-bottom: 4px;
          }

          .value {
            font-size: 14px;
            font-weight: 500;
          }
        }
      }

      .model-pricing {
        display: flex;
        justify-content: space-between;
        margin-bottom: 16px;

        .pricing-item {
          .label {
            display: block;
            font-size: 12px;
            color: #909399;
            margin-bottom: 4px;
          }

          .price {
            font-size: 14px;
            font-weight: 500;
            color: #409eff;
          }
        }
      }

      .model-features {
        .el-tag {
          margin-right: 8px;
        }
      }
    }
  }
}
</style>
