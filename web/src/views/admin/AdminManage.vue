<template>
  <div class="admin-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.adminManage') }}</h2>
      <a-button type="primary" @click="showModal()">
        <PlusOutlined /> {{ $t('common.create') }}
      </a-button>
    </div>

    <a-table :columns="columns" :data-source="admins" :loading="loading" row-key="id">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'role'">
          <a-tag :color="record.role === 'superadmin' ? 'red' : 'blue'">
            {{ record.role === 'superadmin' ? $t('admin.superAdmin') : $t('admin.admin') }}
          </a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-switch :checked="record.status === 'active'" @change="(v: boolean) => handleStatusChange(record.id, v)" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="showModal(record)">{{ $t('common.edit') }}</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="editingAdmin ? $t('common.edit') : $t('common.create')"
      @ok="handleModalOk"
      :confirm-loading="modalLoading"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item :label="$t('auth.username')" name="username" :rules="[{ required: true, message: $t('admin.usernameRequired') }]">
          <a-input v-model:value="formState.username" :disabled="!!editingAdmin" />
        </a-form-item>
        <a-form-item :label="$t('auth.email')" name="email">
          <a-input v-model:value="formState.email" />
        </a-form-item>
        <a-form-item v-if="!editingAdmin" :label="$t('auth.password')" name="password" :rules="[{ required: true, message: $t('admin.passwordRequired') }]">
          <a-input-password v-model:value="formState.password" />
        </a-form-item>
        <a-form-item :label="$t('user.nickname')" name="nickname">
          <a-input v-model:value="formState.nickname" />
        </a-form-item>
        <a-form-item :label="$t('admin.role')" name="role">
          <a-select v-model:value="formState.role">
            <a-select-option value="admin">{{ $t('admin.admin') }}</a-select-option>
            <a-select-option value="superadmin">{{ $t('admin.superAdmin') }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getAdminList, createAdmin, updateAdmin, updateAdminStatus } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const admins = ref<any[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingAdmin = ref<any>(null)

const formState = reactive({
  username: '',
  password: '',
  email: '',
  nickname: '',
  role: 'admin',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('auth.username'), dataIndex: 'username', key: 'username', width: 120 },
  { title: t('auth.email'), dataIndex: 'email', key: 'email', width: 180 },
  { title: t('admin.role'), key: 'role', width: 120 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('admin.lastLoginTime'), dataIndex: 'last_login_at', key: 'last_login_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function fetchAdmins() {
  loading.value = true
  try {
    const res = await getAdminList()
    admins.value = res.data?.list || res.data || []
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function showModal(record?: any) {
  editingAdmin.value = record || null
  if (record) {
    Object.assign(formState, { username: record.username, email: record.email || '', nickname: record.nickname, role: record.role, password: '' })
  } else {
    Object.assign(formState, { username: '', password: '', email: '', nickname: '', role: 'admin' })
  }
  modalVisible.value = true
}

async function handleModalOk() {
  modalLoading.value = true
  try {
    if (editingAdmin.value) {
      await updateAdmin(editingAdmin.value.id, formState)
    } else {
      await createAdmin(formState)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchAdmins()
  } catch { /* handled */ } finally {
    modalLoading.value = false
  }
}

async function handleStatusChange(id: number, active: boolean) {
  try {
    await updateAdminStatus(id, { status: active ? 'active' : 'disabled' })
    message.success(t('common.success'))
    fetchAdmins()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchAdmins()
})
</script>

<style scoped lang="less">
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  h2 { margin: 0; }
}
</style>
