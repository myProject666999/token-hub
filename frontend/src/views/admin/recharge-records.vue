<template>
  <div class="recharge-records-page">
    <h2 class="page-title">充值记录</h2>
    <el-card>
      <el-alert title="充值记录管理功能" type="info" show-icon>
        <template #default>
          <p>在此页面可以查看所有用户的充值记录。</p>
        </template>
      </el-alert>
      <div class="placeholder" style="margin-top: 20px;">
        <el-empty description="充值记录列表">
          <template #image>
            <el-icon size="80"><Money /></el-icon>
          </template>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.recharge-records-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }
}
</style>
