<template>
  <div class="api-keys-page">
    <div class="page-header">
      <h2 class="page-title">API密钥管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        新建密钥
      </el-button>
    </div>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><InfoFilled /></el-icon>
          <span>重要提示</span>
        </div>
      </template>
      <ul class="info-list">
        <li>API 密钥用于调用 Token Hub 的所有 API 接口，请妥善保管</li>
        <li>请勿将 API 密钥提交到公开的代码仓库或分享给他人</li>
        <li>如果密钥泄露，请立即删除并创建新的密钥</li>
        <li>您可以创建多个 API 密钥，用于不同的应用场景</li>
      </ul>
    </el-card>

    <el-card class="list-card">
      <el-table :data="apiKeys" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="key" label="密钥">
          <template #default="{ row }">
            <el-input
              :value="showKey === row.id ? row.key : 'sk-******'"
              :readonly="true"
              class="key-input"
            >
              <template #append>
                <el-button @click="toggleShowKey(row.id)">
                  <el-icon><View /></el-icon>
                </el-button>
                <el-button @click="copyKey(row.key)">
                  <el-icon><CopyDocument /></el-icon>
                </el-button>
              </template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_used" label="最后使用" width="180">
          <template #default="{ row }">
            {{ row.last_used || '从未使用' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="toggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button
              type="danger"
              link
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="showCreateDialog"
      title="新建API密钥"
      width="400px"
    >
      <el-form :model="createForm" :rules="createRules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input
            v-model="createForm.name"
            placeholder="请输入密钥名称，用于标识用途"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAPIKeys, createAPIKey, deleteAPIKey, updateAPIKeyStatus } from '@/api/user'

const loading = ref(false)
const creating = ref(false)
const apiKeys = ref([])
const showKey = ref(null)
const showCreateDialog = ref(false)

const createForm = reactive({
  name: ''
})

const createRules = {
  name: [
    { required: true, message: '请输入密钥名称', trigger: 'blur' }
  ]
}

const loadAPIKeys = async () => {
  loading.value = true
  try {
    const res = await getAPIKeys()
    apiKeys.value = res.data
  } catch (e) {
    console.error('Failed to load API keys:', e)
  } finally {
    loading.value = false
  }
}

const toggleShowKey = (id) => {
  showKey.value = showKey.value === id ? null : id
}

const copyKey = async (key) => {
  try {
    await navigator.clipboard.writeText(key)
    ElMessage.success('已复制到剪贴板')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}

const handleCreate = async () => {
  creating.value = true
  try {
    await createAPIKey(createForm.name)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.name = ''
    loadAPIKeys()
  } catch (e) {
    console.error('Failed to create API key:', e)
  } finally {
    creating.value = false
  }
}

const toggleStatus = async (row) => {
  try {
    const newStatus = row.status === 1 ? 0 : 1
    await updateAPIKeyStatus(row.id, newStatus)
    ElMessage.success(newStatus === 1 ? '已启用' : '已禁用')
    loadAPIKeys()
  } catch (e) {
    console.error('Failed to update status:', e)
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该 API 密钥吗？删除后无法恢复', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAPIKey(row.id)
    ElMessage.success('删除成功')
    loadAPIKeys()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Failed to delete API key:', e)
    }
  }
}

onMounted(() => {
  loadAPIKeys()
})
</script>

<style scoped lang="scss">
.api-keys-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .page-title {
    font-size: 24px;
    margin: 0;
  }

  .info-card {
    margin-bottom: 20px;

    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: bold;
      color: #e6a23c;
    }

    .info-list {
      color: #606266;
      line-height: 2;
      padding-left: 20px;
    }
  }

  .list-card {
    .key-input {
      .el-input__wrapper {
        padding-right: 0;
      }
    }
  }
}
</style>
