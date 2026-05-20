<template>
  <div class="backup-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.backupManage') }}</h2>
      <a-button type="primary" @click="handleCreateBackup" :loading="creating">
        <PlusOutlined /> {{ $t('admin.manualBackup') }}
      </a-button>
    </div>

    <!-- Auto Backup Config -->
    <a-card size="small" class="mb-16">
      <a-row :gutter="16" align="middle">
        <a-col :span="8">
          <span>{{ $t('admin.autoBackup') }}:</span>
          <a-switch v-model:checked="autoBackupConfig.enabled" class="ml-8" />
        </a-col>
        <a-col :span="12">
          <span>{{ $t('admin.backupFrequency') }}:</span>
          <a-select v-model:value="autoBackupConfig.frequency" style="width: 150px; margin-left: 8px;">
            <a-select-option value="daily">{{ $t('admin.daily') }}</a-select-option>
            <a-select-option value="weekly">{{ $t('admin.weekly') }}</a-select-option>
            <a-select-option value="monthly">{{ $t('admin.monthly') }}</a-select-option>
          </a-select>
        </a-col>
        <a-col :span="4">
          <a-button type="primary" size="small" @click="handleSaveAutoBackup">{{ $t('common.save') }}</a-button>
        </a-col>
      </a-row>
    </a-card>

    <a-table :columns="columns" :data-source="backups" :loading="loading" row-key="id">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'size'">
          {{ formatSize(record.size) }}
        </template>
        <template v-if="column.key === 'type'">
          <a-tag :color="record.type === 'manual' ? 'blue' : 'green'">
            {{ record.type === 'manual' ? $t('admin.manual') : $t('admin.auto') }}
          </a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="handleDownload(record.id, record.filename)">
              <DownloadOutlined /> {{ $t('common.download') }}
            </a-button>
            <a-popconfirm :title="`${$t('common.confirm')}${$t('common.delete')}?`" @confirm="handleDelete(record.id)">
              <a-button type="link" danger size="small">{{ $t('common.delete') }}</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getBackupList, createBackup, downloadBackup, deleteBackup, updateAutoBackupConfig } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const backups = ref<any[]>([])
const loading = ref(false)
const creating = ref(false)

const autoBackupConfig = reactive({
  enabled: false,
  frequency: 'daily',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('admin.filename'), dataIndex: 'filename', key: 'filename' },
  { title: t('admin.fileSize'), key: 'size', width: 120 },
  { title: t('admin.backupType'), key: 'type', width: 100 },
  { title: t('article.publishTime'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 180 },
]

function formatSize(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function fetchBackups() {
  loading.value = true
  try {
    const res = await getBackupList()
    const data = res.data || {}
    backups.value = data.list || data || []
    if (data.auto_backup_config) {
      autoBackupConfig.enabled = data.auto_backup_config.enabled ?? false
      autoBackupConfig.frequency = data.auto_backup_config.frequency ?? 'daily'
    }
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

async function handleCreateBackup() {
  creating.value = true
  try {
    await createBackup()
    message.success(t('common.success'))
    fetchBackups()
  } catch { /* handled */ } finally {
    creating.value = false
  }
}

async function handleDownload(id: number, filename: string) {
  try {
    const res = await downloadBackup(id)
    const url = window.URL.createObjectURL(new Blob([res as any]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename || `backup-${id}.sql`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteBackup(id)
    message.success(t('common.success'))
    fetchBackups()
  } catch { /* handled */ }
}

async function handleSaveAutoBackup() {
  try {
    await updateAutoBackupConfig(autoBackupConfig)
    message.success(t('common.success'))
  } catch { /* handled */ }
}

onMounted(() => {
  fetchBackups()
})
</script>

<style scoped lang="less">
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  h2 { margin: 0; }
}

.ml-8 {
  margin-left: 8px;
}
</style>
