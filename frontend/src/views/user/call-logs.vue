<template>
  <div class="call-logs-page">
    <h2 class="page-title">调用日志</h2>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable @change="loadLogs">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadLogs">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card">
      <el-table :data="callLogs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="request_id" label="请求ID" min-width="200" />
        <el-table-column label="服务商" width="120">
          <template #default="{ row }">
            {{ row.provider?.name }}
          </template>
        </el-table-column>
        <el-table-column label="模型" width="150">
          <template #default="{ row }">
            {{ row.model?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="input_tokens" label="输入Tokens" width="120" />
        <el-table-column prop="output_tokens" label="输出Tokens" width="120" />
        <el-table-column prop="total_tokens" label="总Tokens" width="100" />
        <el-table-column prop="points_consumed" label="消耗积分" width="100">
          <template #default="{ row }">
            <span class="points-text">{{ row.points_consumed }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="调用时间" width="180" />
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="loadLogs"
        @current-change="loadLogs"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUserCallLogs } from '@/api/user'

const loading = ref(false)
const callLogs = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = reactive({
  status: ''
})

const loadLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filters.status) {
      params.status = filters.status
    }
    const res = await getUserCallLogs(params)
    callLogs.value = res.data
    total.value = res.total
  } catch (e) {
    console.error('Failed to load call logs:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped lang="scss">
.call-logs-page {
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
  }

  .filter-card {
    margin-bottom: 20px;
  }

  .list-card {
    .pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }

    .points-text {
      color: #e6a23c;
      font-weight: bold;
    }
  }
}
</style>
