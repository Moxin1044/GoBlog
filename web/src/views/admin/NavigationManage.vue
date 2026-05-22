<template>
  <div class="navigation-manage">
    <div class="page-header">
      <h2>{{ $t('admin.navigationManage') }}</h2>
      <a-button type="primary" @click="openModal()"><template #icon><PlusOutlined /></template> {{ $t('common.create') }}</a-button>
    </div>

    <a-card class="tree-card">
      <draggable v-model="flatNavigations" item-key="id" @end="handleDragEnd" handle=".drag-handle" animation="300">
        <template #item="{ element }">
          <div class="nav-item" :style="{ paddingLeft: (getLevel(element) * 24 + 16) + 'px' }">
            <div class="drag-handle">
              <MenuOutlined />
            </div>
            <div class="nav-info">
              <span class="nav-name">{{ element.name }}</span>
              <a-tag v-if="element.type === 'category'" color="blue">{{ $t('article.category') }}</a-tag>
              <a-tag v-else-if="element.type === 'link'" color="orange">{{ $t('navigation.link') }}</a-tag>
              <a-tag v-else color="green">{{ $t('navigation.custom') }}</a-tag>
              <a-tag v-if="!element.enabled" color="default" style="margin-left: 8px">{{ $t('common.disabled') }}</a-tag>
            </div>
            <div class="nav-actions">
              <a-button type="text" size="small" @click="openModal(element)">
                <template #icon><EditOutlined /></template>
              </a-button>
              <a-button type="text" size="small" @click="openModal(null, element.id)">
                <template #icon><FolderAddOutlined /></template>
              </a-button>
              <a-popconfirm
                :title="$t('common.deleteConfirm')"
                @confirm="handleDelete(element.id)"
              >
                <a-button type="text" size="small" danger>
                  <template #icon><DeleteOutlined /></template>
                </a-button>
              </a-popconfirm>
            </div>
          </div>
        </template>
      </draggable>
    </a-card>

    <!-- 编辑模态框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingId ? $t('common.edit') : $t('common.create')"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="600px"
    >
      <a-form
        :model="formState"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
      >
        <a-form-item :label="$t('navigation.name')" required>
          <a-input v-model:value="formState.name" :placeholder="$t('navigation.name')" />
        </a-form-item>
        <a-form-item :label="$t('navigation.nameEn')">
          <a-input v-model:value="formState.name_en" :placeholder="$t('navigation.nameEn')" />
        </a-form-item>
        <a-form-item :label="$t('navigation.type')">
          <a-select v-model:value="formState.type" @change="handleTypeChange">
            <a-select-option value="custom">{{ $t('navigation.custom') }}</a-select-option>
            <a-select-option value="category">{{ $t('navigation.categoryPage') }}</a-select-option>
            <a-select-option value="link">{{ $t('navigation.customLink') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="formState.type === 'category'" :label="$t('navigation.selectCategory')" required>
          <a-select v-model:value="formState.category_id" :placeholder="$t('navigation.selectCategory')">
            <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="formState.type === 'link'" :label="$t('navigation.link')" required>
          <a-input v-model:value="formState.link" :placeholder="$t('navigation.link')" />
        </a-form-item>
        <a-form-item :label="$t('navigation.newTab')">
          <a-switch v-model:checked="formState.new_tab" />
        </a-form-item>
        <a-form-item :label="$t('common.enabled')">
          <a-switch v-model:checked="formState.enabled" />
        </a-form-item>
        <a-form-item :label="$t('navigation.sort')">
          <a-input-number v-model:value="formState.sort" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { PlusOutlined, EditOutlined, DeleteOutlined, MenuOutlined, FolderAddOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import draggable from 'vuedraggable'
import { adminGetNavigations, createNavigation, updateNavigation, deleteNavigation, updateNavigationSort } from '@/api/admin'
import { adminGetCategories } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const navigations = ref<any[]>([])
const categories = ref<any[]>([])
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const parentId = ref<number | null>(null)

const formState = reactive({
  name: '',
  name_en: '',
  type: 'custom',
  link: '',
  category_id: null,
  new_tab: false,
  enabled: true,
  sort: 0,
  parent_id: 0
})

const flatNavigations = computed(() => {
  const result: any[] = []
  const flatten = (items: any[], level: number = 0) => {
    items.forEach(item => {
      result.push({ ...item, _level: level })
      if (item.children && item.children.length > 0) {
        flatten(item.children, level + 1)
      }
    })
  }
  flatten(navigations.value)
  return result
})

const fetchNavigations = async () => {
  try {
    const res = await adminGetNavigations()
    navigations.value = res.data || []
  } catch {
    message.error(t('common.error'))
  }
}

const fetchCategories = async () => {
  try {
    const res = await adminGetCategories()
    categories.value = res.data || []
  } catch {
    message.error(t('common.error'))
  }
}

const getLevel = (item: any) => item._level || 0

const openModal = (item: any = null, pId: number | null = null) => {
  editingId.value = item ? item.id : null
  parentId.value = pId
  if (item) {
    formState.name = item.name
    formState.name_en = item.name_en
    formState.type = item.type
    formState.link = item.link
    formState.category_id = item.category_id
    formState.new_tab = item.new_tab
    formState.enabled = item.enabled
    formState.sort = item.sort
    formState.parent_id = item.parent_id
  } else {
    formState.name = ''
    formState.name_en = ''
    formState.type = 'custom'
    formState.link = ''
    formState.category_id = null
    formState.new_tab = false
    formState.enabled = true
    formState.sort = 0
    formState.parent_id = pId || 0
  }
  modalVisible.value = true
}

const handleTypeChange = () => {
  formState.link = ''
  formState.category_id = null
}

const handleSubmit = async () => {
  if (!formState.name) {
    message.error(t('navigation.nameRequired'))
    return
  }
  if (formState.type === 'category' && !formState.category_id) {
    message.error(t('navigation.categoryRequired'))
    return
  }
  if (formState.type === 'link' && !formState.link) {
    message.error(t('navigation.linkRequired'))
    return
  }
  try {
    if (editingId.value) {
      await updateNavigation(editingId.value, formState)
      message.success(t('common.success'))
    } else {
      await createNavigation(formState)
      message.success(t('common.success'))
    }
    modalVisible.value = false
    fetchNavigations()
  } catch {
    message.error(t('common.error'))
  }
}

const handleCancel = () => {
  modalVisible.value = false
}

const handleDelete = async (id: number) => {
  try {
    await deleteNavigation(id)
    message.success(t('common.success'))
    fetchNavigations()
  } catch (err: any) {
    message.error(err.message || t('common.error'))
  }
}

const handleDragEnd = async () => {
  const sortData = flatNavigations.value.map((item, index) => ({
    id: item.id,
    sort: index
  }))
  try {
    await updateNavigationSort(sortData)
    message.success(t('common.success'))
    fetchNavigations()
  } catch {
    message.error(t('common.error'))
  }
}

onMounted(() => {
  fetchNavigations()
  fetchCategories()
})
</script>

<style scoped lang="less">
.navigation-manage {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .tree-card {
    .nav-item {
      display: flex;
      align-items: center;
      padding: 12px 16px;
      border-bottom: 1px solid #f0f0f0;
      transition: background-color 0.2s;

      &:hover {
        background-color: #fafafa;
      }

      &:last-child {
        border-bottom: none;
      }

      .drag-handle {
        cursor: grab;
        padding: 0 8px;
        color: #999;

        &:active {
          cursor: grabbing;
        }
      }

      .nav-info {
        flex: 1;
        display: flex;
        align-items: center;
        gap: 8px;

        .nav-name {
          font-weight: 500;
        }
      }

      .nav-actions {
        display: flex;
        gap: 4px;
      }
    }
  }
}
</style>
