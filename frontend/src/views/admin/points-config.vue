<template>
  <div class="points-config-page">
    <h2 class="page-title">积分配置</h2>

    <el-card class="config-card">
      <template #header>
        <span>积分汇率配置</span>
      </template>
      <el-form label-width="120px">
        <el-form-item label="当前汇率">
          <div class="current-rate">
            <span>1 元 = </span>
            <span class="rate-value">{{ currentRate }}</span>
            <span> 积分</span>
          </div>
        </el-form-item>
        <el-form-item label="设置新汇率">
          <el-input-number
            v-model="newRate"
            :min="1"
            :max="10000"
            :precision="0"
            placeholder="请输入新的积分汇率"
          />
          <span class="unit">积分/元</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleUpdate">
            保存配置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <span>说明</span>
      </template>
      <ul class="info-list">
        <li>积分汇率表示 1 元人民币可以兑换多少积分</li>
        <li>用户充值时，系统会自动按照当前汇率兑换积分</li>
        <li>修改汇率只会影响后续的充值，已充值的积分不会受到影响</li>
        <li>不同的模型调用消耗的积分数量可以在模型管理中设置</li>
      </ul>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPointsRate, setPointsRate } from '@/api/payment'

const loading = ref(false)
const currentRate = ref(100)
const newRate = ref(100)

const loadRate = async () => {
  try {
    const res = await getPointsRate()
    currentRate.value = res.data.rate
    newRate.value = res.data.rate
  } catch (e) {
    console.error('Failed to load rate:', e)
  }
}

const handleUpdate = async () => {
  loading.value = true
  try {
    await setPointsRate(newRate.value)
    ElMessage.success('保存成功')
    currentRate.value = newRate.value
  } catch (e) {
    console.error('Failed to update rate:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRate()
})
</script>

<style scoped lang="scss">
.points-config-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }

  .config-card {
    margin-bottom: 20px;
  }

  .current-rate {
    font-size: 18px;

    .rate-value {
      color: #409eff;
      font-weight: bold;
      font-size: 24px;
    }
  }

  .unit {
    margin-left: 10px;
    color: #909399;
  }

  .info-list {
    color: #606266;
    line-height: 2;
    padding-left: 20px;
  }
}
</style>
