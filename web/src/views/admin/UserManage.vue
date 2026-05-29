<template>
  <div class="user-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.userManage') }}</h2>
      <a-button type="primary" @click="showCreateModal">{{ $t('admin.addUser') }}</a-button>
    </div>

    <div class="filter-bar mb-16">
      <a-space wrap>
        <a-input-search
          v-model:value="usernameSearch"
          :placeholder="$t('auth.username')"
          @search="fetchUsers"
          style="width: 200px"
        />
        <a-input-search
          v-model:value="emailSearch"
          :placeholder="$t('auth.email')"
          @search="fetchUsers"
          style="width: 200px"
        />
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="users"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'avatar'">
          <a-avatar :src="record.avatar">{{ record.nickname?.charAt(0) || 'U' }}</a-avatar>
        </template>
        <template v-if="column.key === 'status'">
          <a-switch :checked="record.status === 'active'" @change="(v: boolean) => handleStatusChange(record.id, v)" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="showDetail(record)">{{ $t('admin.viewDetail') }}</a-button>
            <a-popconfirm :title="$t('admin.confirmResetPwd')" @confirm="handleResetPassword(record.id)">
              <a-button type="link" size="small">{{ $t('admin.resetPassword') }}</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- User Detail Modal -->
    <a-modal
      v-model:open="detailVisible"
      :title="$t('admin.userDetail')"
      :footer="null"
      width="600px"
    >
      <a-descriptions bordered :column="2" v-if="detailUser">
        <a-descriptions-item :label="$t('auth.username')">{{ detailUser.username }}</a-descriptions-item>
        <a-descriptions-item :label="$t('user.nickname')">{{ detailUser.nickname }}</a-descriptions-item>
        <a-descriptions-item :label="$t('auth.email')">{{ detailUser.email || '-' }}</a-descriptions-item>
        <a-descriptions-item :label="$t('auth.phone')">{{ detailUser.phone || '-' }}</a-descriptions-item>
        <a-descriptions-item :label="$t('common.status')">
          <a-tag :color="detailUser.status === 'active' ? 'green' : 'red'">
            {{ detailUser.status === 'active' ? $t('common.enable') : $t('common.disable') }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item :label="$t('admin.registerTime')">{{ detailUser.created_at || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- Create User Modal -->
    <a-modal
      v-model:open="createVisible"
      :title="$t('admin.addUser')"
      @ok="handleCreateUser"
      :confirm-loading="createLoading"
    >
      <a-form :model="createForm" layout="vertical">
        <a-form-item :label="$t('auth.username')" required>
          <a-input v-model:value="createForm.username" />
        </a-form-item>
        <a-form-item :label="$t('auth.email')" required>
          <a-input v-model:value="createForm.email" />
        </a-form-item>
        <a-form-item :label="$t('auth.password')" required>
          <a-input-password v-model:value="createForm.password" />
        </a-form-item>
        <a-form-item :label="$t('user.nickname')">
          <a-input v-model:value="createForm.nickname" />
        </a-form-item>
        <a-form-item :label="$t('auth.phone')">
          <a-input v-model:value="createForm.phone" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { adminGetUsers, adminCreateUser, adminGetUser, updateUserStatus, resetUserPassword } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const users = ref<any[]>([])
const loading = ref(false)
const usernameSearch = ref('')
const emailSearch = ref('')
const detailVisible = ref(false)
const detailUser = ref<any>(null)
const createVisible = ref(false)
const createLoading = ref(false)

const createForm = reactive({
  username: '',
  email: '',
  password: '',
  nickname: '',
  phone: '',
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('auth.username'), dataIndex: 'username', key: 'username', width: 120 },
  { title: t('auth.email'), dataIndex: 'email', key: 'email', width: 180 },
  { title: t('auth.phone'), dataIndex: 'phone', key: 'phone', width: 130 },
  { title: t('user.nickname'), dataIndex: 'nickname', key: 'nickname', width: 120 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('admin.registerTime'), dataIndex: 'created_at', key: 'created_at', width: 160 },
  { title: t('common.actions'), key: 'actions', width: 180 },
]

async function fetchUsers() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.current, page_size: pagination.pageSize }
    if (usernameSearch.value) params.username = usernameSearch.value
    if (emailSearch.value) params.email = emailSearch.value
    const res = await adminGetUsers(params)
    users.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchUsers()
}

function showCreateModal() {
  createForm.username = ''
  createForm.email = ''
  createForm.password = ''
  createForm.nickname = ''
  createForm.phone = ''
  createVisible.value = true
}

async function handleCreateUser() {
  if (!createForm.username || !createForm.email || !createForm.password) {
    message.warning(t('admin.fillRequired'))
    return
  }
  if (createForm.password.length < 6) {
    message.warning(t('auth.passwordMin'))
    return
  }
  createLoading.value = true
  try {
    await adminCreateUser(createForm)
    message.success(t('common.success'))
    createVisible.value = false
    fetchUsers()
  } catch { /* handled */ } finally {
    createLoading.value = false
  }
}

async function showDetail(record: any) {
  try {
    const res = await adminGetUser(record.id)
    detailUser.value = res.data || record
    detailVisible.value = true
  } catch {
    detailUser.value = record
    detailVisible.value = true
  }
}

async function handleStatusChange(id: number, active: boolean) {
  try {
    await updateUserStatus(id, { status: active ? 'active' : 'disabled' })
    message.success(t('common.success'))
    fetchUsers()
  } catch { /* handled */ }
}

async function handleResetPassword(id: number) {
  try {
    await resetUserPassword(id)
    message.success(t('admin.passwordResetSuccess'))
  } catch { /* handled */ }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="less">
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  h2 {
    margin: 0;
  }
}
</style>
