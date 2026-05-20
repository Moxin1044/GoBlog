<template>
  <div class="ai-model-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.aiModelManage') }}</h2>
      <a-button type="primary" @click="showModal()">
        <PlusOutlined /> {{ $t('common.create') }}
      </a-button>
    </div>

    <!-- Global AI Controls -->
    <a-card size="small" class="mb-16">
      <a-row :gutter="16" align="middle">
        <a-col :span="8">
          <span>{{ $t('admin.aiChatSwitch') }}:</span>
          <a-switch v-model:checked="globalConfig.ai_enabled" @change="handleGlobalConfigChange" class="ml-8" />
        </a-col>
        <a-col :span="16">
          <span>{{ $t('admin.dailyChatLimit') }}:</span>
          <a-input-number v-model:value="globalConfig.daily_chat_limit" :min="0" :max="9999" style="width: 120px; margin-left: 8px;" />
          <a-button type="link" size="small" @click="handleGlobalConfigChange">{{ $t('common.save') }}</a-button>
        </a-col>
      </a-row>
    </a-card>

    <a-table :columns="columns" :data-source="models" :loading="loading" row-key="id">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'provider'">
          {{ providerText(record.provider) }}
        </template>
        <template v-if="column.key === 'status'">
          <a-switch :checked="record.status === 'active'" @change="(v: boolean) => handleStatusChange(record.id, v)" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="showModal(record)">{{ $t('common.edit') }}</a-button>
            <a-popconfirm :title="`${$t('common.confirm')}${$t('common.delete')}?`" @confirm="handleDelete(record.id)">
              <a-button type="link" danger size="small">{{ $t('common.delete') }}</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="editingModel ? $t('common.edit') : $t('common.create')"
      @ok="handleModalOk"
      :confirm-loading="modalLoading"
      width="700px"
    >
      <a-form :model="formState" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('admin.modelName')" name="name" :rules="[{ required: true }]">
              <a-input v-model:value="formState.name" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('admin.provider')" name="provider" :rules="[{ required: true }]">
              <a-select v-model:value="formState.provider">
                <a-select-option value="openai">OpenAI</a-select-option>
                <a-select-option value="qwen">{{ $t('admin.qwen') }}</a-select-option>
                <a-select-option value="spark">{{ $t('admin.spark') }}</a-select-option>
                <a-select-option value="doubao">{{ $t('admin.doubao') }}</a-select-option>
                <a-select-option value="local">{{ $t('admin.localModel') }}</a-select-option>
                <a-select-option value="custom">{{ $t('admin.custom') }}</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="16">
            <a-form-item label="API URL" name="api_url" :rules="[{ required: true }]">
              <a-input v-model:value="formState.api_url" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('admin.requestType')" name="request_type">
              <a-select v-model:value="formState.request_type">
                <a-select-option value="POST">POST</a-select-option>
                <a-select-option value="GET">GET</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item :label="$t('user.apiToken')" name="api_token">
          <a-input-password v-model:value="formState.api_token" />
        </a-form-item>
        <a-form-item :label="$t('admin.requestHeaders')" name="request_headers">
          <a-textarea v-model:value="formState.request_headers" :rows="3" placeholder='{"Content-Type": "application/json"}' />
        </a-form-item>
        <a-form-item :label="$t('admin.requestTemplate')" name="request_template">
          <a-textarea v-model:value="formState.request_template" :rows="4" :placeholder='$t("admin.requestTemplatePlaceholder")' />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('ai.contextLength')" name="max_context">
              <a-input-number v-model:value="formState.max_context" :min="1" :max="128" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('user.temperature')" name="temperature">
              <a-slider v-model:value="formState.temperature" :min="0" :max="2" :step="0.1" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getAIModelList, createAIModel, updateAIModel, updateAIModelStatus, deleteAIModel, updateAIModelGlobalConfig } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const models = ref<any[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const modalLoading = ref(false)
const editingModel = ref<any>(null)

const globalConfig = reactive({
  ai_enabled: true,
  daily_chat_limit: 10,
})

const formState = reactive({
  name: '',
  provider: 'openai',
  api_url: '',
  api_token: '',
  request_type: 'POST',
  request_headers: '',
  request_template: '',
  temperature: 0.7,
  max_context: 10,
})

const providerMap: Record<string, string> = {
  openai: 'OpenAI',
  qwen: t('admin.qwen'),
  spark: t('admin.spark'),
  doubao: t('admin.doubao'),
  local: t('admin.localModel'),
  custom: t('admin.custom'),
}

function providerText(provider: string) {
  return providerMap[provider] || provider
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('admin.modelName'), dataIndex: 'name', key: 'name', width: 150 },
  { title: t('admin.provider'), key: 'provider', width: 120 },
  { title: 'API URL', dataIndex: 'api_url', key: 'api_url', ellipsis: true },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('common.actions'), key: 'actions', width: 150 },
]

async function fetchModels() {
  loading.value = true
  try {
    const res = await getAIModelList()
    const data = res.data || {}
    models.value = data.list || data || []
    if (data.global_config) {
      globalConfig.ai_enabled = data.global_config.ai_enabled ?? true
      globalConfig.daily_chat_limit = data.global_config.daily_chat_limit ?? 10
    }
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function showModal(record?: any) {
  editingModel.value = record || null
  if (record) {
    Object.assign(formState, {
      name: record.name,
      provider: record.provider || 'openai',
      api_url: record.api_url,
      api_token: record.api_token || '',
      request_type: record.request_type || 'POST',
      request_headers: record.request_headers || '',
      request_template: record.request_template || '',
      temperature: record.temperature ?? 0.7,
      max_context: record.max_context ?? 10,
    })
  } else {
    Object.assign(formState, {
      name: '', provider: 'openai', api_url: '', api_token: '',
      request_type: 'POST', request_headers: '', request_template: '',
      temperature: 0.7, max_context: 10,
    })
  }
  modalVisible.value = true
}

async function handleModalOk() {
  modalLoading.value = true
  try {
    if (editingModel.value) {
      await updateAIModel(editingModel.value.id, formState)
    } else {
      await createAIModel(formState)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchModels()
  } catch { /* handled */ } finally {
    modalLoading.value = false
  }
}

async function handleStatusChange(id: number, active: boolean) {
  try {
    await updateAIModelStatus(id, { status: active ? 'active' : 'disabled' })
    message.success(t('common.success'))
    fetchModels()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteAIModel(id)
    message.success(t('common.success'))
    fetchModels()
  } catch { /* handled */ }
}

async function handleGlobalConfigChange() {
  try {
    await updateAIModelGlobalConfig(globalConfig)
    message.success(t('common.success'))
  } catch { /* handled */ }
}

onMounted(() => {
  fetchModels()
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
