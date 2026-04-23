<template>
  <div class="payment-methods-page">
    <h2 class="page-title">支付方式管理</h2>
    <el-card>
      <el-alert title="支付方式管理功能" type="info" show-icon>
        <template #default>
          <p>在此页面可以管理支付宝、微信支付等支付方式的配置。</p>
        </template>
      </el-alert>
      <div class="placeholder" style="margin-top: 20px;">
        <el-empty description="支付方式列表">
          <template #image>
            <el-icon size="80"><Wallet /></el-icon>
          </template>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped lang="scss">
.payment-methods-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }
}
</style>
