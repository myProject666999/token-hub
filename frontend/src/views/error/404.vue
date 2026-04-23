<template>
  <div class="error-page">
    <div class="error-content">
      <h1 class="error-code">404</h1>
      <h2 class="error-message">页面不存在</h2>
      <p class="error-desc">抱歉，您访问的页面不存在或已被移除</p>
      <el-button type="primary" size="large" @click="$router.push('/')">
        <el-icon><House /></el-icon>
        返回首页
      </el-button>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.error-content {
  text-align: center;
  padding: 40px;

  .error-code {
    font-size: 120px;
    font-weight: bold;
    color: #409eff;
    margin: 0 0 20px;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
  }

  .error-message {
    font-size: 32px;
    color: #303133;
    margin: 0 0 16px;
  }

  .error-desc {
    font-size: 16px;
    color: #909399;
    margin: 0 0 30px;
  }
}
</style>
