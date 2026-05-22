<template>
  <div class="category-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.categoryManage') }}</h2>
      <a-button type="primary" @click="showCreateModal">
        <PlusOutlined /> {{ $t('common.create') }}
      </a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="categories"
      :loading="loading"
      row-key="id"
      :scroll="{ x: 600 }"
      :pagination="false"
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
        <a-form-item :label="$t('category.name')" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item :label="$t('category.nameEn')">
          <a-input v-model:value="formState.name_en" />
        </a-form-item>
        <a-form-item :label="$t('category.sort')">
          <a-input-number v-model:value="formState.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { adminGetCategories, createCategory, updateCategory, deleteCategory } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const categories = ref<any[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const formState = reactive({
  name: '',
  name_en: '',
  sort: 0,
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('category.name'), dataIndex: 'name', key: 'name' },
  { title: t('category.nameEn'), dataIndex: 'name_en', key: 'name_en' },
  { title: t('category.sort'), dataIndex: 'sort', key: 'sort', width: 100 },
  { title: t('common.createdAt'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 150, fixed: 'right' as const },
]

async function fetchCategories() {
  loading.value = true
  try {
    const res = await adminGetCategories()
    categories.value = res.data || []
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function showCreateModal() {
  isEdit.value = false
  editingId.value = null
  formState.name = ''
  formState.name_en = ''
  formState.sort = 0
  modalVisible.value = true
}

function showEditModal(record: any) {
  isEdit.value = true
  editingId.value = record.id
  formState.name = record.name
  formState.name_en = record.name_en || ''
  formState.sort = record.sort || 0
  modalVisible.value = true
}

function handleModalCancel() {
  modalVisible.value = false
}

async function handleModalOk() {
  if (!formState.name) {
    message.warning(t('category.nameRequired'))
    return
  }

  try {
    if (isEdit.value && editingId.value) {
      await updateCategory(editingId.value, formState)
    } else {
      await createCategory(formState)
    }
    message.success(t('common.success'))
    modalVisible.value = false
    fetchCategories()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteCategory(id)
    message.success(t('common.success'))
    fetchCategories()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchCategories()
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
