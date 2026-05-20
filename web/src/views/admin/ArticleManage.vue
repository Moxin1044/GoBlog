<template>
  <div class="article-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.articleManage') }}</h2>
      <a-button type="primary" @click="$router.push('/admin/article/create')">
        <PlusOutlined /> {{ $t('common.create') }}
      </a-button>
    </div>

    <div class="filter-bar mb-16">
      <a-space wrap>
        <a-input-search
          v-model:value="searchText"
          :placeholder="$t('common.search')"
          @search="fetchArticles"
          style="width: 250px"
        />
        <a-select
          v-model:value="categoryFilter"
          :placeholder="$t('article.category')"
          allow-clear
          style="width: 150px"
          @change="fetchArticles"
        >
          <a-select-option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</a-select-option>
        </a-select>
        <a-select
          v-model:value="statusFilter"
          :placeholder="$t('common.status')"
          allow-clear
          style="width: 120px"
          @change="fetchArticles"
        >
          <a-select-option value="published">{{ $t('article.published') }}</a-select-option>
          <a-select-option value="draft">{{ $t('article.draft') }}</a-select-option>
          <a-select-option value="offline">{{ $t('article.offline') }}</a-select-option>
        </a-select>
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="articles"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'tags'">
          <a-tag v-for="tag in (record.tags || [])" :key="tag.id" color="blue">{{ tag.name }}</a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="$router.push(`/admin/article/edit/${record.id}`)">
              {{ $t('common.edit') }}
            </a-button>
            <a-button
              v-if="record.status === 'published'"
              type="link"
              size="small"
              @click="handleToggleStatus(record.id, 'offline')"
            >
              {{ $t('article.offline') }}
            </a-button>
            <a-button
              v-if="record.status === 'offline'"
              type="link"
              size="small"
              @click="handleToggleStatus(record.id, 'published')"
            >
              {{ $t('article.published') }}
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
import { PlusOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { adminGetArticleList, deleteArticle, updateArticleStatus, adminGetCategories } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const articles = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)
const searchText = ref('')
const categoryFilter = ref<number | undefined>(undefined)
const statusFilter = ref<string | undefined>(undefined)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => t('admin.totalCount', { total }),
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('article.title'), dataIndex: 'title', key: 'title', ellipsis: true },
  { title: t('article.category'), dataIndex: 'category_name', key: 'category', width: 120 },
  { title: t('article.tags'), key: 'tags', width: 180 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('article.publishTime'), dataIndex: 'published_at', key: 'published_at', width: 160 },
  { title: t('article.viewCount'), dataIndex: 'view_count', key: 'view_count', width: 90 },
  { title: t('common.actions'), key: 'actions', width: 200, fixed: 'right' as const },
]

function statusColor(status: string) {
  const map: Record<string, string> = { published: 'green', draft: 'orange', offline: 'red' }
  return map[status] || 'default'
}

function statusText(status: string) {
  const map: Record<string, string> = { published: t('article.published'), draft: t('article.draft'), offline: t('article.offline') }
  return map[status] || status
}

async function fetchCategories() {
  try {
    const res = await adminGetCategories()
    categories.value = res.data?.list || res.data || []
  } catch { /* handled */ }
}

async function fetchArticles() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: pagination.current,
      page_size: pagination.pageSize,
    }
    if (searchText.value) params.search = searchText.value
    if (statusFilter.value) params.status = statusFilter.value
    if (categoryFilter.value) params.category_id = categoryFilter.value
    const res = await adminGetArticleList(params)
    articles.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchArticles()
}

async function handleToggleStatus(id: number, status: string) {
  try {
    await updateArticleStatus(id, { status })
    message.success(t('common.success'))
    fetchArticles()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteArticle(id)
    message.success(t('common.success'))
    fetchArticles()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchCategories()
  fetchArticles()
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
