<template>
  <div class="home-page">
    <header class="page-header">
      <div class="logo" @click="$router.push('/')">
        <el-icon><Connection /></el-icon>
        <span>Token Hub</span>
      </div>
      <nav class="nav">
        <el-button type="primary" @click="$router.push('/models')">模型列表</el-button>
        <el-button @click="$router.push('/docs')">接入说明</el-button>
        <template v-if="userStore.isAuthenticated">
          <el-button @click="$router.push('/user/dashboard')">控制台</el-button>
          <el-button @click="handleLogout">退出</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" @click="$router.push('/register')">注册</el-button>
        </template>
      </nav>
    </header>

    <section class="hero">
      <div class="hero-content">
        <h1>AI大模型中转平台</h1>
        <p>一键接入主流大模型，统一API接口，灵活的积分制管理</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" @click="$router.push('/models')">查看模型</el-button>
          <el-button size="large" @click="$router.push('/docs')">快速接入</el-button>
        </div>
      </div>
    </section>

    <section class="features">
      <h2 class="section-title">为什么选择我们</h2>
      <div class="feature-list">
        <div class="feature-item">
          <el-icon class="feature-icon"><Cpu /></el-icon>
          <h3>多模型支持</h3>
          <p>支持OpenAI、Claude、GLM、千问、混元、豆包等主流大模型</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Link /></el-icon>
          <h3>统一API</h3>
          <p>兼容OpenAI API格式，一次对接，多模型切换</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Coin /></el-icon>
          <h3>积分制</h3>
          <p>灵活的积分管理，按使用量计费，透明消费</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Document /></el-icon>
          <h3>详细日志</h3>
          <p>完整的调用记录，随时查看使用情况</p>
        </div>
      </div>
    </section>

    <section class="providers-section">
      <h2 class="section-title">支持的服务商</h2>
      <div class="providers-grid">
        <div class="provider-card" v-for="provider in providers" :key="provider.id">
          <div class="provider-header">
            <h3>{{ provider.name }}</h3>
          </div>
          <p>{{ provider.description }}</p>
          <div class="model-count">
            <el-tag type="primary">{{ provider.models?.length || 0 }} 个模型</el-tag>
          </div>
        </div>
      </div>
    </section>

    <footer class="page-footer">
      <p>&copy; 2024 Token Hub. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getProviders } from '@/api/provider'

const router = useRouter()
const userStore = useUserStore()

const providers = ref([])

const handleLogout = () => {
  userStore.logout()
  router.push('/')
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
.home-page {
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

.hero {
  padding: 150px 40px 100px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  text-align: center;
  color: #fff;

  .hero-content {
    max-width: 800px;
    margin: 0 auto;

    h1 {
      font-size: 48px;
      margin-bottom: 20px;
    }

    p {
      font-size: 18px;
      margin-bottom: 40px;
      opacity: 0.9;
    }

    .hero-actions {
      .el-button {
        margin: 0 10px;
      }
    }
  }
}

.features {
  padding: 80px 40px;
  background-color: #fff;

  .section-title {
    text-align: center;
    font-size: 32px;
    margin-bottom: 60px;
  }

  .feature-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 40px;
    max-width: 1200px;
    margin: 0 auto;

    .feature-item {
      text-align: center;

      .feature-icon {
        font-size: 48px;
        color: #409eff;
        margin-bottom: 20px;
      }

      h3 {
        font-size: 20px;
        margin-bottom: 15px;
      }

      p {
        color: #666;
        line-height: 1.6;
      }
    }
  }
}

.providers-section {
  padding: 80px 40px;
  background-color: #f5f7fa;

  .section-title {
    text-align: center;
    font-size: 32px;
    margin-bottom: 60px;
  }

  .providers-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 30px;
    max-width: 1200px;
    margin: 0 auto;

    .provider-card {
      background: #fff;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      transition: transform 0.3s;

      &:hover {
        transform: translateY(-4px);
      }

      .provider-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 16px;

        h3 {
          margin: 0;
          font-size: 18px;
        }
      }

      p {
        color: #666;
        margin-bottom: 16px;
      }
    }
  }
}

.page-footer {
  padding: 40px;
  background-color: #304156;
  color: #fff;
  text-align: center;
}
</style>
