<template>
  <div class="recharge-page">
    <h2 class="page-title">充值中心</h2>

    <el-card class="points-rate-card">
      <template #header>
        <span>当前汇率</span>
      </template>
      <div class="rate-content">
        <div class="rate-display">
          <span class="rate-value">1 元 = {{ pointsRate }} 积分</span>
        </div>
        <p class="rate-tip">充值金额将按照当前汇率自动兑换为积分</p>
      </div>
    </el-card>

    <el-card class="recharge-card">
      <template #header>
        <span>选择充值金额</span>
      </template>
      <div class="amount-options">
        <div
          v-for="amount in presetAmounts"
          :key="amount"
          class="amount-option"
          :class="{ active: selectedAmount === amount }"
          @click="selectedAmount = amount"
        >
          <div class="amount-value">¥{{ amount }}</div>
          <div class="amount-points">≈ {{ amount * pointsRate }} 积分</div>
        </div>
      </div>

      <el-form class="custom-amount-form" label-position="left">
        <el-form-item label="自定义金额">
          <el-input-number
            v-model="customAmount"
            :min="1"
            :max="10000"
            :precision="0"
            placeholder="请输入充值金额"
          />
          <span class="unit">元</span>
        </el-form-item>
      </el-form>

      <div class="payment-methods">
        <h4>选择支付方式</h4>
        <div class="method-list">
          <div
            v-for="method in paymentMethods"
            :key="method.id"
            class="method-item"
            :class="{ active: selectedMethod === method.id }"
            @click="selectedMethod = method.id"
          >
            <el-icon v-if="method.code === 'alipay'" class="method-icon alipay"><Money /></el-icon>
            <el-icon v-else-if="method.code === 'wechat'" class="method-icon wechat"><ChatDotRound /></el-icon>
            <el-icon v-else class="method-icon"><Wallet /></el-icon>
            <span class="method-name">{{ method.name }}</span>
          </div>
        </div>
      </div>

      <div class="recharge-summary">
        <div class="summary-row">
          <span>充值金额</span>
          <span class="highlight">¥{{ currentAmount }}</span>
        </div>
        <div class="summary-row">
          <span>获得积分</span>
          <span class="highlight points">{{ currentAmount * pointsRate }} 积分</span>
        </div>
      </div>

      <el-button
        type="primary"
        size="large"
        class="recharge-btn"
        :loading="loading"
        @click="handleRecharge"
        :disabled="!selectedMethod || currentAmount < 1"
      >
        立即充值
      </el-button>
    </el-card>

    <el-card class="history-card">
      <template #header>
        <span>充值记录</span>
      </template>
      <el-table :data="rechargeRecords" v-loading="loadingRecords">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount }}
          </template>
        </el-table-column>
        <el-table-column prop="points" label="积分" width="150">
          <template #default="{ row }">
            {{ row.points }} 积分
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="primary"
              link
              @click="simulatePayment(row.order_no)"
            >
              模拟支付
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="loadRecords"
        @current-change="loadRecords"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPointsRate, getPaymentMethods, createRechargeOrder, getUserRechargeRecords, simulatePayment } from '@/api/payment'

const presetAmounts = [10, 50, 100, 200, 500, 1000]
const selectedAmount = ref(null)
const customAmount = ref(0)
const selectedMethod = ref(null)
const paymentMethods = ref([])
const pointsRate = ref(100)
const loading = ref(false)

const rechargeRecords = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loadingRecords = ref(false)

const currentAmount = computed(() => {
  return customAmount.value > 0 ? customAmount.value : (selectedAmount.value || 0)
})

const loadPointsRate = async () => {
  try {
    const res = await getPointsRate()
    pointsRate.value = res.data.rate
  } catch (e) {
    console.error('Failed to load points rate:', e)
  }
}

const loadPaymentMethods = async () => {
  try {
    const res = await getPaymentMethods(false)
    paymentMethods.value = res.data
    if (res.data.length > 0) {
      selectedMethod.value = res.data[0].id
    }
  } catch (e) {
    console.error('Failed to load payment methods:', e)
  }
}

const loadRecords = async () => {
  loadingRecords.value = true
  try {
    const res = await getUserRechargeRecords({
      page: page.value,
      page_size: pageSize.value
    })
    rechargeRecords.value = res.data
    total.value = res.total
  } catch (e) {
    console.error('Failed to load recharge records:', e)
  } finally {
    loadingRecords.value = false
  }
}

const handleRecharge = async () => {
  if (!currentAmount.value || currentAmount.value < 1) {
    ElMessage.warning('请选择或输入充值金额')
    return
  }
  if (!selectedMethod.value) {
    ElMessage.warning('请选择支付方式')
    return
  }

  loading.value = true
  try {
    const res = await createRechargeOrder({
      payment_method_id: selectedMethod.value,
      amount: currentAmount.value
    })
    ElMessage.success('订单创建成功')
    loadRecords()
    
    if (confirm(`订单创建成功！\n订单号：${res.data.order_no}\n是否进行模拟支付？`)) {
      await simulatePayment(res.data.order_no)
      ElMessage.success('模拟支付成功！积分已到账')
      loadRecords()
      loadPointsRate()
    }
  } catch (e) {
    console.error('Failed to create recharge order:', e)
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    pending: 'warning',
    paid: 'success',
    expired: 'info',
    refund: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待支付',
    paid: '已支付',
    expired: '已过期',
    refund: '已退款'
  }
  return texts[status] || status
}

onMounted(() => {
  loadPointsRate()
  loadPaymentMethods()
  loadRecords()
})
</script>

<style scoped lang="scss">
.recharge-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }

  .points-rate-card {
    margin-bottom: 20px;

    .rate-content {
      text-align: center;
    }

    .rate-display {
      .rate-value {
        font-size: 28px;
        font-weight: bold;
        color: #409eff;
      }
    }

    .rate-tip {
      color: #909399;
      margin-top: 8px;
    }
  }

  .recharge-card {
    margin-bottom: 20px;

    .amount-options {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 12px;
      margin-bottom: 24px;

      .amount-option {
        border: 2px solid #e4e7ed;
        border-radius: 8px;
        padding: 16px;
        text-align: center;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: #409eff;
        }

        &.active {
          border-color: #409eff;
          background-color: #ecf5ff;
        }

        .amount-value {
          font-size: 20px;
          font-weight: bold;
          color: #303133;
        }

        .amount-points {
          font-size: 12px;
          color: #909399;
          margin-top: 4px;
        }
      }
    }

    .custom-amount-form {
      margin-bottom: 24px;
      padding-bottom: 24px;
      border-bottom: 1px solid #ebeef5;

      .unit {
        margin-left: 8px;
        color: #909399;
      }
    }

    .payment-methods {
      margin-bottom: 24px;

      h4 {
        margin-bottom: 12px;
        font-size: 14px;
        color: #606266;
      }

      .method-list {
        display: flex;
        gap: 12px;
      }

      .method-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 20px;
        border: 2px solid #e4e7ed;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: #409eff;
        }

        &.active {
          border-color: #409eff;
          background-color: #ecf5ff;
        }

        .method-icon {
          font-size: 24px;
        }

        .alipay {
          color: #1677ff;
        }

        .wechat {
          color: #07c160;
        }

        .method-name {
          font-weight: 500;
        }
      }
    }

    .recharge-summary {
      background: #f5f7fa;
      border-radius: 8px;
      padding: 20px;
      margin-bottom: 24px;

      .summary-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        &:last-child {
          margin-bottom: 0;
        }

        .highlight {
          font-weight: bold;
          font-size: 18px;
        }

        .points {
          color: #409eff;
        }
      }
    }

    .recharge-btn {
      width: 100%;
    }
  }

  .history-card {
    .pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }
}
</style>
