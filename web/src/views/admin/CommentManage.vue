<template>
  <div class="comment-manage">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.commentManage') }}</h2>
    </div>

    <div class="filter-bar mb-16">
      <a-space wrap>
        <a-select v-model:value="articleFilter" style="width: 200px" @change="fetchComments" :placeholder="$t('article.title')" allow-clear show-search :filter-option="filterOption">
          <a-select-option v-for="art in articles" :key="art.id" :value="art.id">{{ art.title }}</a-select-option>
        </a-select>
        <a-select v-model:value="statusFilter" style="width: 120px" @change="fetchComments" :placeholder="$t('common.status')" allow-clear>
          <a-select-option value="pending">{{ $t('comment.pending') }}</a-select-option>
          <a-select-option value="approved">{{ $t('comment.approved') }}</a-select-option>
          <a-select-option value="rejected">{{ $t('comment.rejected') }}</a-select-option>
        </a-select>
        <a-range-picker v-model:value="dateRange" @change="fetchComments" />
        <a-button :disabled="!selectedRowKeys.length" @click="handleBatchReview('approved')" type="primary">
          {{ $t('comment.batchReview') }} - {{ $t('comment.approved') }}
        </a-button>
        <a-button :disabled="!selectedRowKeys.length" @click="handleBatchReview('rejected')" danger>
          {{ $t('comment.batchReview') }} - {{ $t('comment.rejected') }}
        </a-button>
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="comments"
      :loading="loading"
      :pagination="pagination"
      :row-selection="{ selectedRowKeys, onChange: onSelectChange }"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button v-if="record.status === 'pending'" type="link" size="small" @click="handleReview(record.id, 'approved')">
              {{ $t('comment.approved') }}
            </a-button>
            <a-button v-if="record.status === 'pending'" type="link" danger size="small" @click="handleReview(record.id, 'rejected')">
              {{ $t('comment.rejected') }}
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
import { message } from 'ant-design-vue'
import { adminGetComments, reviewComment, deleteComment, batchReviewComments, adminGetArticleList } from '@/api/admin'
import { useI18n } from 'vue-i18n'
import type { Dayjs } from 'dayjs'

const { t } = useI18n()
const comments = ref<any[]>([])
const articles = ref<any[]>([])
const loading = ref(false)
const statusFilter = ref<string | undefined>(undefined)
const articleFilter = ref<number | undefined>(undefined)
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const selectedRowKeys = ref<number[]>([])

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('article.title'), dataIndex: 'article_title', key: 'article_title', width: 150 },
  { title: t('comment.nickname'), dataIndex: 'nickname', key: 'nickname', width: 120 },
  { title: t('comment.content'), dataIndex: 'content', key: 'content', ellipsis: true },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('admin.submitTime'), dataIndex: 'created_at', key: 'created_at', width: 160 },
  { title: t('admin.reviewTime'), dataIndex: 'reviewed_at', key: 'reviewed_at', width: 160 },
  { title: t('common.actions'), key: 'actions', width: 200 },
]

function statusColor(status: string) {
  const map: Record<string, string> = { approved: 'green', pending: 'orange', rejected: 'red' }
  return map[status] || 'default'
}

function statusText(status: string) {
  const map: Record<string, string> = { approved: t('comment.approved'), pending: t('comment.pending'), rejected: t('comment.rejected') }
  return map[status] || status
}

function filterOption(input: string, option: any) {
  return option.children?.[0]?.children?.toLowerCase().includes(input.toLowerCase())
}

async function fetchArticles() {
  try {
    const res = await adminGetArticleList({ page: 1, page_size: 100 })
    articles.value = res.data?.list || []
  } catch { /* handled */ }
}

async function fetchComments() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.current, page_size: pagination.pageSize }
    if (statusFilter.value) params.status = statusFilter.value
    if (articleFilter.value) params.article_id = articleFilter.value
    if (dateRange.value) {
      params.start_date = dateRange.value[0].format('YYYY-MM-DD')
      params.end_date = dateRange.value[1].format('YYYY-MM-DD')
    }
    const res = await adminGetComments(params)
    comments.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchComments()
}

function onSelectChange(keys: number[]) {
  selectedRowKeys.value = keys
}

async function handleReview(id: number, status: string) {
  try {
    await reviewComment(id, { status })
    message.success(t('common.success'))
    fetchComments()
  } catch { /* handled */ }
}

async function handleBatchReview(status: string) {
  try {
    await batchReviewComments({ ids: selectedRowKeys.value, status })
    message.success(t('common.success'))
    selectedRowKeys.value = []
    fetchComments()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await deleteComment(id)
    message.success(t('common.success'))
    fetchComments()
  } catch { /* handled */ }
}

onMounted(() => {
  fetchArticles()
  fetchComments()
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
