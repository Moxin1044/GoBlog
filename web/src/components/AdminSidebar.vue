<template>
  <a-layout-sider
    :collapsed="collapsed"
    :trigger="null"
    collapsible
    width="220"
    class="admin-sidebar"
    :class="{ dark: themeStore.isDark }"
  >
    <div class="sidebar-logo">
      <router-link to="/admin" class="logo-link">
        <img v-if="appStore.siteConfig.siteLogo" :src="appStore.siteConfig.siteLogo" alt="Logo" class="logo-img" />
        <span v-if="!collapsed" class="logo-text">GoBlog Admin</span>
      </router-link>
    </div>
    <a-menu
      mode="inline"
      :selected-keys="selectedKeys"
      :open-keys="openKeys"
      theme="dark"
      @click="handleMenuClick"
    >
      <a-menu-item key="dashboard">
        <DashboardOutlined />
        <span>{{ $t('admin.dashboard') }}</span>
      </a-menu-item>
      <a-menu-item key="articles">
        <FileTextOutlined />
        <span>{{ $t('admin.articleManage') }}</span>
      </a-menu-item>
      <a-menu-item key="categories">
        <FolderOutlined />
        <span>{{ $t('admin.categoryManage') }}</span>
      </a-menu-item>
      <a-menu-item key="tags">
        <TagsOutlined />
        <span>{{ $t('admin.tagManage') }}</span>
      </a-menu-item>
      <a-menu-item key="comments">
        <CommentOutlined />
        <span>{{ $t('admin.commentManage') }}</span>
      </a-menu-item>
      <a-menu-item key="users">
        <UserOutlined />
        <span>{{ $t('admin.userManage') }}</span>
      </a-menu-item>
      <a-menu-item v-if="userStore.isSuperAdmin()" key="admins">
        <TeamOutlined />
        <span>{{ $t('admin.adminManage') }}</span>
      </a-menu-item>
      <a-menu-item v-if="userStore.isSuperAdmin()" key="config">
        <SettingOutlined />
        <span>{{ $t('admin.systemConfig') }}</span>
      </a-menu-item>
      <a-menu-item v-if="userStore.isSuperAdmin()" key="ai-models">
        <RobotOutlined />
        <span>{{ $t('admin.aiModelManage') }}</span>
      </a-menu-item>
      <a-menu-item v-if="userStore.isSuperAdmin()" key="backups">
        <DatabaseOutlined />
        <span>{{ $t('admin.backupManage') }}</span>
      </a-menu-item>
      <a-menu-item key="logs">
        <FileSearchOutlined />
        <span>{{ $t('admin.operationLog') }}</span>
      </a-menu-item>
    </a-menu>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  DashboardOutlined,
  FileTextOutlined,
  FolderOutlined,
  TagsOutlined,
  CommentOutlined,
  UserOutlined,
  TeamOutlined,
  SettingOutlined,
  RobotOutlined,
  DatabaseOutlined,
  FileSearchOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { useThemeStore } from '@/stores/theme'

defineProps<{
  collapsed: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()
const themeStore = useThemeStore()

const openKeys = ref<string[]>([])

const selectedKeys = computed(() => {
  const path = route.path
  if (path === '/admin') return ['dashboard']
  if (path.includes('article')) return ['articles']
  if (path.includes('categor')) return ['categories']
  if (path.includes('tag')) return ['tags']
  if (path.includes('comment')) return ['comments']
  if (path.includes('user')) return ['users']
  if (path.includes('admin-mgmt') || path.includes('admins')) return ['admins']
  if (path.includes('config')) return ['config']
  if (path.includes('ai-model')) return ['ai-models']
  if (path.includes('backup')) return ['backups']
  if (path.includes('log')) return ['logs']
  return ['dashboard']
})

function handleMenuClick({ key }: { key: string }) {
  if (key === 'dashboard') router.push('/admin')
  else router.push(`/admin/${key}`)
}
</script>

<style scoped lang="less">
.admin-sidebar {
  overflow: auto;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 10;
}

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
}

.logo-img {
  height: 32px;
}
</style>
