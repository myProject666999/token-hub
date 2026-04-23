<template>
  <div class="call-logs-page">
    <h2 class="page-title">调用日志</h2>
    <el-card>
      <el-alert title="调用日志管理功能" type="info" show-icon>
        <template #default>
          <p>在此页面可以查看所有用户的API调用日志，包括调用的模型、消耗的积分等信息。</p>
        </template>
      </el-alert>
      <div class="placeholder" style="margin-top: 20px;">
        <el-empty description="调用日志列表">
          <template #image>
            <el-icon size="80"><Document /></el-icon>
          </template>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.call-logs-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }
}
</style>
