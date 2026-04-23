<template>
  <div class="models-page">
    <h2 class="page-title">模型管理</h2>
    <el-card>
      <el-alert title="模型管理功能" type="info" show-icon>
        <template #default>
          <p>在此页面可以管理所有大模型，包括添加、编辑、删除模型，设置积分消耗等。</p>
        </template>
      </el-alert>
      <div class="placeholder" style="margin-top: 20px;">
        <el-empty description="模型列表">
          <template #image>
            <el-icon size="80"><Cpu /></el-icon>
          </template>
          <template #description>
            <span>请参考服务商管理页面实现模型管理功能</span>
          </template>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.models-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }
}
</style>
