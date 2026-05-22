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
          <a-tag :color="providerColor(record.provider)">{{ providerText(record.provider) }}</a-tag>
        </template>
        <template v-if="column.key === 'models'">
          <template v-if="parseModels(record.models).length">
            <a-tag v-for="m in parseModels(record.models).slice(0, 5)" :key="m" size="small" class="model-tag">{{ m }}</a-tag>
            <a-tag v-if="parseModels(record.models).length > 5" size="small">+{{ parseModels(record.models).length - 5 }}</a-tag>
          </template>
          <span v-else class="text-secondary">-</span>
        </template>
        <template v-if="column.key === 'status'">
          <a-switch :checked="record.enabled" @change="(v: boolean) => handleStatusChange(record.id, v)" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="showModal(record)">{{ $t('common.edit') }}</a-button>
            <a-popconfirm :title="$t('common.deleteConfirm')" @confirm="handleDelete(record.id)">
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
      width="600px"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item :label="$t('admin.providerName')" name="name" :rules="[{ required: true }]">
          <a-input v-model:value="formState.name" :placeholder="$t('admin.providerNamePlaceholder')" />
        </a-form-item>
        <a-form-item :label="$t('admin.provider')" name="provider" :rules="[{ required: true }]">
          <a-select v-model:value="formState.provider" :placeholder="$t('admin.selectProvider')">
            <a-select-option value="openai">OpenAI</a-select-option>
            <a-select-option value="qwen">{{ $t('admin.qwen') }}</a-select-option>
            <a-select-option value="spark">{{ $t('admin.spark') }}</a-select-option>
            <a-select-option value="doubao">{{ $t('admin.doubao') }}</a-select-option>
            <a-select-option value="deepseek">DeepSeek</a-select-option>
            <a-select-option value="local">{{ $t('admin.localModel') }}</a-select-option>
            <a-select-option value="custom">{{ $t('admin.custom') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="API URL" name="api_url" :rules="[{ required: true }]">
          <a-input v-model:value="formState.api_url" placeholder="https://api.example.com/v1" />
        </a-form-item>
        <a-form-item :label="$t('admin.modelTags')" name="models">
          <div class="model-tags-editor">
            <div class="tags-display mb-8">
              <a-tag
                v-for="(tag, index) in modelTags"
                :key="index"
                closable
                @close="removeModelTag(index)"
                color="blue"
              >
                {{ tag }}
              </a-tag>
            </div>
            <div class="tag-input-row">
              <a-input
                v-model:value="newModelTag"
                :placeholder="$t('admin.modelTagPlaceholder')"
                @pressEnter="addModelTag"
                style="flex: 1"
              />
              <a-button type="primary" size="small" @click="addModelTag" style="margin-left: 8px;">
                {{ $t('common.create') }}
              </a-button>
            </div>
          </div>
        </a-form-item>
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
const modelTags = ref<string[]>([])
const newModelTag = ref('')

const globalConfig = reactive({
  ai_enabled: true,
  daily_chat_limit: 10,
})

const formState = reactive({
  name: '',
  provider: 'openai',
  api_url: '',
  models: '',
})

const providerMap: Record<string, string> = {
  openai: 'OpenAI',
  qwen: t('admin.qwen'),
  spark: t('admin.spark'),
  doubao: t('admin.doubao'),
  deepseek: 'DeepSeek',
  local: t('admin.localModel'),
  custom: t('admin.custom'),
}

const providerColorMap: Record<string, string> = {
  openai: 'green',
  qwen: 'blue',
  spark: 'purple',
  doubao: 'orange',
  deepseek: 'cyan',
  local: 'default',
  custom: 'default',
}

function providerText(provider: string) {
  return providerMap[provider] || provider
}

function providerColor(provider: string) {
  return providerColorMap[provider] || 'default'
}

function parseModels(modelsStr: string): string[] {
  if (!modelsStr) return []
  try {
    return JSON.parse(modelsStr)
  } catch {
    return []
  }
}

function addModelTag() {
  const tag = newModelTag.value.trim()
  if (!tag) return
  if (modelTags.value.includes(tag)) {
    message.warning(t('admin.modelTagExists'))
    return
  }
  modelTags.value.push(tag)
  newModelTag.value = ''
}

function removeModelTag(index: number) {
  modelTags.value.splice(index, 1)
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('admin.providerName'), dataIndex: 'name', key: 'name', width: 150 },
  { title: t('admin.provider'), key: 'provider', width: 120 },
  { title: 'API URL', dataIndex: 'api_url', key: 'api_url', ellipsis: true },
  { title: t('admin.modelTags'), key: 'models', width: 250 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('common.actions'), key: 'actions', width: 150 },
]

async function fetchModels() {
  loading.value = true
  try {
    const res = await getAIModelList()
    const data = res.data || {}
    models.value = Array.isArray(data) ? data : (data.list || [])
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
    formState.name = record.name
    formState.provider = record.provider || 'openai'
    formState.api_url = record.api_url
    formState.models = record.models || ''
    modelTags.value = parseModels(record.models || '')
  } else {
    formState.name = ''
    formState.provider = 'openai'
    formState.api_url = ''
    formState.models = ''
    modelTags.value = []
  }
  newModelTag.value = ''
  modalVisible.value = true
}

async function handleModalOk() {
  if (!formState.name) {
    message.warning(t('admin.providerNameRequired'))
    return
  }
  if (!formState.api_url) {
    message.warning(t('admin.apiUrlRequired'))
    return
  }

  modalLoading.value = true
  try {
    const data = {
      ...formState,
      models: JSON.stringify(modelTags.value),
    }
    if (editingModel.value) {
      await updateAIModel(editingModel.value.id, data)
    } else {
      await createAIModel(data)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchModels()
  } catch { /* handled */ } finally {
    modalLoading.value = false
  }
}

async function handleStatusChange(id: number, enabled: boolean) {
  try {
    await updateAIModelStatus(id, { enabled })
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

.model-tags-editor {
  .tags-display {
    min-height: 32px;
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .tag-input-row {
    display: flex;
    align-items: center;
  }
}

.model-tag {
  font-size: 12px;
}

.text-secondary {
  color: var(--text-secondary, #999);
}
</style>
