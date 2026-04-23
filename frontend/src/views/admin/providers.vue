<template>
  <div class="providers-page">
    <div class="page-header">
      <h2 class="page-title">服务商管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        新增服务商
      </el-button>
    </div>

    <el-card class="list-card">
      <el-table :data="providers" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="code" label="编码" width="120" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="website" label="官网" min-width="150">
          <template #default="{ row }">
            <el-link v-if="row.website" :href="row.website" target="_blank">
              {{ row.website }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
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
        @size-change="loadProviders"
        @current-change="loadProviders"
      />
    </el-card>

    <el-dialog
      v-model="showCreateDialog"
      :title="isEdit ? '编辑服务商' : '新增服务商'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入服务商名称" />
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入服务商编码（英文）" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="2"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item label="官网" prop="website">
          <el-input v-model="formData.website" placeholder="请输入官网地址" />
        </el-form-item>
        <el-form-item label="API地址" prop="api_endpoint">
          <el-input v-model="formData.api_endpoint" placeholder="请输入API端点地址" />
        </el-form-item>
        <el-form-item label="API密钥" prop="api_key">
          <el-input v-model="formData.api_key" placeholder="请输入API密钥" show-password />
        </el-form-item>
        <el-form-item label="API密钥2" prop="api_secret">
          <el-input v-model="formData.api_secret" placeholder="请输入API密钥（如有）" show-password />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getProviderList, createProvider, updateProvider, deleteProvider, getProvider } from '@/api/provider'

const loading = ref(false)
const submitting = ref(false)
const providers = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const showCreateDialog = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref(null)

const formData = reactive({
  name: '',
  code: '',
  description: '',
  website: '',
  api_endpoint: '',
  api_key: '',
  api_secret: '',
  status: 1,
  sort: 0
})

const formRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }],
  api_endpoint: [{ required: true, message: '请输入API地址', trigger: 'blur' }]
}

const loadProviders = async () => {
  loading.value = true
  try {
    const res = await getProviderList({
      page: page.value,
      page_size: pageSize.value
    })
    providers.value = res.data
    total.value = res.total
  } catch (e) {
    console.error('Failed to load providers:', e)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.website = ''
  formData.api_endpoint = ''
  formData.api_key = ''
  formData.api_secret = ''
  formData.status = 1
  formData.sort = 0
}

const handleEdit = async (row) => {
  isEdit.value = true
  editId.value = row.id
  try {
    const res = await getProvider(row.id)
    Object.assign(formData, res.data)
    showCreateDialog.value = true
  } catch (e) {
    console.error('Failed to load provider:', e)
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该服务商吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteProvider(row.id)
    ElMessage.success('删除成功')
    loadProviders()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Failed to delete:', e)
    }
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateProvider(editId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await createProvider(formData)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    loadProviders()
  } catch (e) {
    console.error('Failed to submit:', e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadProviders()
})

showCreateDialog.value = false
watch(showCreateDialog, (val) => {
  if (!val) {
    resetForm()
    isEdit.value = false
    editId.value = null
  }
})
</script>

<style scoped lang="scss">
.providers-page {
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

  .list-card {
    .pagination {
      margin-top: 20px;
      justify-content: flex-end;
    }
  }
}
</style>
