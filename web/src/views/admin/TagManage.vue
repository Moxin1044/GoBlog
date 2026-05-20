<template>
  <div class="tag-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.tagManage') }}</h2>
      <a-button type="primary" @click="showCreateModal">
        <PlusOutlined /> {{ $t('common.create') }}
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="tags"
      :loading="loading"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="showEditModal(record)">
              {{ $t('common.edit') }}
            </a-button>
            <a-popconfirm :title="`${$t('common.confirm')}${$t('common.delete')}?`" @confirm="handleDelete(record.id)">
              <a-button type="link" danger size="small">{{ $t('common.delete') }}</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? $t('common.edit') : $t('common.create')"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item :label="$t('tag.name')" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item :label="$t('tag.nameEn')">
          <a-input v-model:value="formState.nameEn" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { adminGetTags, createTag, updateTag, deleteTag } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const tags = ref<any[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const formState = reactive({
  name: '',
  nameEn: '',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('tag.name'), dataIndex: 'name', key: 'name' },
  { title: t('tag.nameEn'), dataIndex: 'name_en', key: 'name_en' },
  { title: t('common.createdAt'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 150, fixed: 'right' as const },
]

async function fetchTags() {
  loading.value = true
  try {
    const res = await adminGetTags()
    tags.value = res.data || []
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function showCreateModal() {
  isEdit.value = false
  editingId.value = null
  formState.name = ''
  formState.nameEn = ''
  modalVisible.value = true
}

function showEditModal(record: any) {
  isEdit.value = true
  editingId.value = record.id
  formState.name = record.name
  formState.nameEn = record.name_en || ''
  modalVisible.value = true
}

function handleModalCancel() {
  modalVisible.value = false
}

async function handleModalOk() {
  if (!formState.name) {
    message.warning(t('tag.nameRequired'))
    return
  }

  try {
    if (isEdit.value && editingId.value) {
      await updateTag(editingId.value, formState)
    } else {
      await createTag(formState)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchTags()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteTag(id)
    message.success(t('common.success'))
    fetchTags()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchTags()
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
