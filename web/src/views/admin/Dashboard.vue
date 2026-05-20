<template>
  <div class="dashboard-page">
    <!-- Top Statistics Cards -->
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic :title="$t('admin.totalVisits')" :value="stats.totalVisits || 0">
            <template #prefix><EyeOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic :title="$t('admin.totalArticles')" :value="stats.totalArticles || 0">
            <template #prefix><FileTextOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic :title="$t('admin.totalWords')" :value="stats.totalWords || 0">
            <template #prefix><EditOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card>
          <a-statistic :title="$t('comment.pending')" :value="stats.pendingComments || 0">
            <template #prefix><CommentOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Time Dimension Switch -->
    <div class="dimension-switch mt-16">
      <a-radio-group v-model:value="timeDimension" button-style="solid" @change="onDimensionChange">
        <a-radio-button value="day">{{ $t('admin.day') }}</a-radio-button>
        <a-radio-button value="week">{{ $t('admin.week') }}</a-radio-button>
        <a-radio-button value="month">{{ $t('admin.month') }}</a-radio-button>
      </a-radio-group>
    </div>

    <!-- Server Monitoring + Visit Map -->
    <a-row :gutter="[16, 16]" class="mt-16">
      <a-col :xs="24" :lg="12">
        <a-card :title="$t('admin.realtime')">
          <template #extra>
            <a-radio-group v-model:value="monitorDimension" size="small" @change="onMonitorDimensionChange">
              <a-radio-button value="realtime">{{ $t('admin.realtime') }}</a-radio-button>
              <a-radio-button value="day">{{ $t('admin.day') }}</a-radio-button>
              <a-radio-button value="week">{{ $t('admin.week') }}</a-radio-button>
              <a-radio-button value="month">{{ $t('admin.month') }}</a-radio-button>
            </a-radio-group>
          </template>
          <!-- Server Resource Bars -->
          <div class="monitor-grid">
            <div class="monitor-item">
              <div class="monitor-label">{{ $t('admin.cpuUsage') }}</div>
              <a-progress :percent="parseFloat((monitor.cpu || 0).toFixed(2))" :stroke-color="getProgressColor(monitor.cpu)" />
              <div class="monitor-value">{{ (monitor.cpu || 0).toFixed(2) }}%</div>
            </div>
            <div class="monitor-item">
              <div class="monitor-label">{{ $t('admin.memoryUsage') }}</div>
              <a-progress :percent="parseFloat((monitor.memoryPercent || 0).toFixed(2))" :stroke-color="getProgressColor(monitor.memoryPercent)" />
              <div class="monitor-value">{{ formatBytes(monitor.memoryUsed) }} / {{ formatBytes(monitor.memoryTotal) }} ({{ (monitor.memoryPercent || 0).toFixed(2) }}%)</div>
            </div>
            <div class="monitor-item">
              <div class="monitor-label">{{ $t('admin.diskUsage') }}</div>
              <a-progress :percent="parseFloat((monitor.diskPercent || 0).toFixed(2))" :stroke-color="getProgressColor(monitor.diskPercent)" />
              <div class="monitor-value">{{ formatBytes(monitor.diskUsed) }} / {{ formatBytes(monitor.diskTotal) }} ({{ (monitor.diskPercent || 0).toFixed(2) }}%)</div>
            </div>
            <div class="monitor-item">
              <div class="monitor-label">{{ $t('admin.networkUsage') }}</div>
              <div class="network-stats">
                <span><ArrowUpOutlined /> {{ formatBytes(monitor.networkUpload || 0) }}/s</span>
                <span><ArrowDownOutlined /> {{ formatBytes(monitor.networkDownload || 0) }}/s</span>
              </div>
            </div>
          </div>
          <!-- Server Monitor Chart -->
          <div ref="monitorChartRef" style="height: 300px; margin-top: 16px;"></div>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="12">
        <a-card :title="$t('admin.visitMap')">
          <template #extra>
            <a-radio-group v-model:value="mapViewMode" size="small">
              <a-radio-button value="2d">2D</a-radio-button>
              <a-radio-button value="3d">3D</a-radio-button>
            </a-radio-group>
          </template>
          <div ref="mapRef" style="height: 400px;"></div>
          <!-- Realtime Access Log -->
          <div class="access-log mt-16">
            <div class="log-title">{{ $t('admin.accessLog') }}</div>
            <a-table
              :columns="accessLogColumns"
              :data-source="accessLogs"
              :pagination="false"
              size="small"
              :scroll="{ y: 200 }"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'device'">
                  <span>
                    <DesktopOutlined v-if="record.device === 'desktop'" />
                    <MobileOutlined v-else-if="record.device === 'mobile'" />
                    <TabletOutlined v-else />
                    {{ record.device }}
                  </span>
                </template>
              </template>
            </a-table>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import {
  EyeOutlined, FileTextOutlined, EditOutlined, CommentOutlined,
  ArrowUpOutlined, ArrowDownOutlined, DesktopOutlined, MobileOutlined, TabletOutlined,
} from '@ant-design/icons-vue'
import { getDashboard, getServerMonitor, getVisitMapData } from '@/api/admin'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Stats
const stats = reactive({
  totalVisits: 0,
  totalArticles: 0,
  totalWords: 0,
  pendingComments: 0,
})

// Time dimension
const timeDimension = ref('day')

// Monitor
const monitor = reactive({
  cpu: 0,
  memoryPercent: 0,
  memoryUsed: 0,
  memoryTotal: 0,
  diskPercent: 0,
  diskUsed: 0,
  diskTotal: 0,
  networkUpload: 0,
  networkDownload: 0,
})
const monitorDimension = ref('realtime')
const monitorChartRef = ref<HTMLElement>()
let monitorChart: any = null
let monitorTimer: ReturnType<typeof setInterval> | null = null

// Map
const mapRef = ref<HTMLElement>()
const mapViewMode = ref('2d')
let mapChart: any = null

// Access logs
const accessLogs = ref<any[]>([])
const accessLogColumns = [
  { title: 'IP', dataIndex: 'ip', key: 'ip', width: 120 },
  { title: t('admin.location'), dataIndex: 'location', key: 'location', width: 120 },
  { title: t('admin.accessTime'), dataIndex: 'time', key: 'time', width: 160 },
  { title: t('admin.device'), key: 'device', width: 100 },
]

function getProgressColor(value: number) {
  if (value > 80) return '#ff4d4f'
  if (value > 60) return '#faad14'
  return '#1890ff'
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function fetchDashboard() {
  try {
    const res = await getDashboard()
    const data = res.data || {}
    stats.totalVisits = data.total_views || 0
    stats.totalArticles = data.article_count || 0
    stats.totalWords = data.total_words || 0
    stats.pendingComments = data.comment_count || 0
  } catch { /* handled */ }
}

async function fetchMonitor() {
  try {
    const res = await getServerMonitor()
    const data = res.data || {}
    const cpu = data.cpu || {}
    const memory = data.memory || {}
    const disk = data.disk || {}
    const network = data.network || {}

    monitor.cpu = Array.isArray(cpu.percent) ? cpu.percent[0] || 0 : (cpu.percent || 0)
    monitor.memoryPercent = memory.used_percent || 0
    monitor.memoryUsed = memory.used || 0
    monitor.memoryTotal = memory.total || 0
    monitor.diskPercent = disk.used_percent || 0
    monitor.diskUsed = disk.used || 0
    monitor.diskTotal = disk.total || 0
    monitor.networkUpload = network.bytes_sent || 0
    monitor.networkDownload = network.bytes_recv || 0
  } catch { /* handled */ }
}

async function initMonitorChart() {
  if (!monitorChartRef.value) return
  const echarts = await import('echarts')
  monitorChart = echarts.init(monitorChartRef.value)
  updateMonitorChart()
}

function updateMonitorChart() {
  if (!monitorChart) return
  const now = new Date()
  const timeData: string[] = []
  const cpuData: number[] = []
  const memData: number[] = []
  const netUpData: number[] = []
  const netDownData: number[] = []

  for (let i = 29; i >= 0; i--) {
    const t = new Date(now.getTime() - i * 2000)
    timeData.push(`${t.getHours().toString().padStart(2, '0')}:${t.getMinutes().toString().padStart(2, '0')}:${t.getSeconds().toString().padStart(2, '0')}`)
    cpuData.push(Math.max(0, monitor.cpu + (Math.random() - 0.5) * 10))
    memData.push(Math.max(0, monitor.memoryPercent + (Math.random() - 0.5) * 5))
    netUpData.push(Math.max(0, (monitor.networkUpload || 0) + (Math.random() - 0.5) * 1000))
    netDownData.push(Math.max(0, (monitor.networkDownload || 0) + (Math.random() - 0.5) * 2000))
  }

  monitorChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['CPU', t('admin.memoryUsage'), 'Upload', 'Download'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', boundaryGap: false, data: timeData },
    yAxis: { type: 'value' },
    series: [
      { name: 'CPU', type: 'line', smooth: true, data: cpuData, itemStyle: { color: '#1890ff' } },
      { name: t('admin.memoryUsage'), type: 'line', smooth: true, data: memData, itemStyle: { color: '#52c41a' } },
      { name: 'Upload', type: 'line', smooth: true, data: netUpData, itemStyle: { color: '#faad14' } },
      { name: 'Download', type: 'line', smooth: true, data: netDownData, itemStyle: { color: '#722ed1' } },
    ],
  })
}

async function initMapChart() {
  if (!mapRef.value) return
  const echarts = await import('echarts')
  mapChart = echarts.init(mapRef.value)

  try {
    const res = await getVisitMapData()
    const mapData = res.data || { points: [], logs: [] }
    renderMap(mapData.points || [])
    accessLogs.value = (mapData.logs || []).slice(0, 10)
  } catch {
    renderMap([])
  }
}

function renderMap(points: any[]) {
  if (!mapChart) return
  const isDark = document.body.classList.contains('dark')

  // Scatter points for visitor locations
  const scatterData = points.map((p: any) => ({
    name: p.name || p.location,
    value: [p.lng || p.longitude || 0, p.lat || p.latitude || 0, p.count || 1],
  }))

  // Flylines from Beijing to visitor locations
  const beijingCoord = [116.46, 39.92]
  const linesData = scatterData.map((p) => ({
    coords: [beijingCoord, [p.value[0], p.value[1]]],
  }))

  mapChart.setOption({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'item' },
    geo: {
      map: 'world',
      roam: true,
      label: { show: false },
      itemStyle: {
        areaColor: isDark ? '#1a1a2e' : '#e0e0e0',
        borderColor: isDark ? '#333' : '#999',
      },
      emphasis: {
        itemStyle: { areaColor: isDark ? '#2a2a4e' : '#c0c0c0' },
      },
    },
    series: [
      {
        name: t('admin.visitMap'),
        type: 'effectScatter',
        coordinateSystem: 'geo',
        data: scatterData,
        symbolSize: (val: number[]) => Math.min(Math.max(val[2], 6), 20),
        encode: { value: 2 },
        showEffectOn: 'render',
        rippleEffect: { brushType: 'stroke', scale: 3, period: 4 },
        label: { show: false },
        itemStyle: { color: '#1890ff', shadowBlur: 10, shadowColor: '#1890ff' },
      },
      {
        type: 'lines',
        coordinateSystem: 'geo',
        data: linesData,
        lineStyle: { color: '#1890ff', width: 1, opacity: 0.4, curveness: 0.3 },
        effect: {
          show: true,
          period: 6,
          trailLength: 0.7,
          color: '#fff',
          symbolSize: 3,
        },
      },
    ],
  }, true)
}

function onDimensionChange() {
  fetchDashboard()
}

function onMonitorDimensionChange() {
  updateMonitorChart()
}

watch(mapViewMode, () => {
  nextTick(() => {
    if (mapChart) {
      mapChart.dispose()
      mapChart = null
    }
    initMapChart()
  })
})

onMounted(async () => {
  await fetchDashboard()
  await fetchMonitor()
  await nextTick()
  initMonitorChart()
  initMapChart()

  // Auto-refresh monitor data
  monitorTimer = setInterval(async () => {
    await fetchMonitor()
    updateMonitorChart()
  }, 5000)
})

onBeforeUnmount(() => {
  if (monitorTimer) clearInterval(monitorTimer)
  if (monitorChart) { monitorChart.dispose(); monitorChart = null }
  if (mapChart) { mapChart.dispose(); mapChart = null }
})
</script>

<style scoped lang="less">
.dashboard-page {
  .dimension-switch {
    display: flex;
    justify-content: flex-end;
  }

  .monitor-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }

  .monitor-item {
    .monitor-label {
      font-size: 14px;
      color: var(--text-secondary, #666);
      margin-bottom: 8px;
    }

    .monitor-value {
      font-size: 12px;
      color: var(--text-secondary, #999);
      margin-top: 4px;
    }
  }

  .network-stats {
    display: flex;
    gap: 16px;
    font-size: 14px;
    margin-top: 8px;

    span {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .access-log {
    .log-title {
      font-size: 14px;
      font-weight: 500;
      margin-bottom: 8px;
    }
  }
}
</style>
