<template>
  <div class="operation-log">
    <div class="page-header mb-16">
      <h2>{{ $t('admin.operationLog') }}</h2>
    </div>

    <div class="filter-bar mb-16">
      <a-space wrap>
        <a-input-search
          v-model:value="searchText"
          :placeholder="$t('admin.operatorSearch')"
          @search="fetchLogs"
          style="width: 200px"
        />
        <a-range-picker v-model:value="dateRange" @change="fetchLogs" />
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="logs"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'result'">
          <a-tag :color="isSuccess(record.result) ? 'green' : 'orange'">
            {{ record.result }}
          </a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getOperationLogs } from '@/api/admin'
import { useI18n } from 'vue-i18n'
import type { Dayjs } from 'dayjs'

const { t } = useI18n()
const logs = ref<any[]>([])
const loading = ref(false)
const searchText = ref('')
const dateRange = ref<[Dayjs, Dayjs] | null>(null)

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: t('admin.operator'), dataIndex: 'username', key: 'username', width: 120 },
  { title: t('admin.actionContent'), dataIndex: 'action', key: 'action' },
  { title: t('admin.actionTarget'), dataIndex: 'target', key: 'target', width: 150 },
  { title: t('admin.actionResult'), key: 'result', width: 100 },
  { title: 'IP', dataIndex: 'ip', key: 'ip', width: 140 },
  { title: t('admin.actionTime'), dataIndex: 'created_at', key: 'created_at', width: 180 },
]

function isSuccess(result: string) {
  return result === '成功' || result === 'success'
}

async function fetchLogs() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: pagination.current,
      page_size: pagination.pageSize,
    }
    if (searchText.value) params.search = searchText.value
    if (dateRange.value) {
      params.start_date = dateRange.value[0].format('YYYY-MM-DD')
      params.end_date = dateRange.value[1].format('YYYY-MM-DD')
    }
    const res = await getOperationLogs(params)
    logs.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchLogs()
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped lang="less">
.page-header {
  h2 { margin: 0; }
}
</style>
