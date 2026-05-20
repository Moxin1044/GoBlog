<template>
  <a-layout-header class="admin-header">
    <div class="header-left">
      <a-button type="text" @click="$emit('toggle')">
        <MenuUnfoldOutlined v-if="collapsed" />
        <MenuFoldOutlined v-else />
      </a-button>
      <a-breadcrumb class="header-breadcrumb">
        <a-breadcrumb-item>
          <router-link to="/admin">{{ $t('nav.admin') }}</router-link>
        </a-breadcrumb-item>
        <a-breadcrumb-item>{{ currentPageTitle }}</a-breadcrumb-item>
      </a-breadcrumb>
    </div>
    <div class="header-right">
      <a-dropdown :trigger="['click']">
        <a-button type="text" size="small" class="lang-btn">
          <GlobalOutlined />
          <span class="lang-text">{{ currentLangLabel }}</span>
        </a-button>
        <template #overlay>
          <a-menu @click="handleLangChange">
            <a-menu-item key="zh-CN">中文</a-menu-item>
            <a-menu-item key="en-US">English</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
      <ThemeToggle />
      <a-dropdown>
        <a-space class="user-info" :style="{ cursor: 'pointer' }">
          <a-avatar>{{ userStore.username.charAt(0).toUpperCase() }}</a-avatar>
          <span>{{ userStore.username }}</span>
        </a-space>
        <template #overlay>
          <a-menu>
            <a-menu-item key="profile" @click="$router.push('/profile')">{{ $t('nav.profile') }}</a-menu-item>
            <a-menu-item key="home" @click="$router.push('/')">{{ $t('nav.home') }}</a-menu-item>
            <a-menu-divider />
            <a-menu-item key="logout" @click="handleLogout">{{ $t('nav.logout') }}</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { MenuUnfoldOutlined, MenuFoldOutlined, GlobalOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import ThemeToggle from './ThemeToggle.vue'

defineProps<{
  collapsed: boolean
}>()

defineEmits<{
  (e: 'toggle'): void
}>()

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const userStore = useUserStore()

const currentLangLabel = computed(() => {
  return locale.value === 'zh-CN' ? '中文' : 'EN'
})

const pageKeyMap: Record<string, string> = {
  '/admin': 'admin.dashboard',
  '/admin/articles': 'admin.articleManage',
  '/admin/comments': 'admin.commentManage',
  '/admin/users': 'admin.userManage',
  '/admin/admins': 'admin.adminManage',
  '/admin/config': 'admin.systemConfig',
  '/admin/ai-models': 'admin.aiModelManage',
  '/admin/backups': 'admin.backupManage',
  '/admin/logs': 'admin.operationLog',
}

const currentPageTitle = computed(() => {
  const path = route.path
  if (path.includes('article/create') || path.includes('article/edit')) return t('common.create')
  const key = pageKeyMap[path] || 'admin.dashboard'
  return t(key)
})

function handleLangChange({ key }: { key: string }) {
  locale.value = key
  localStorage.setItem('locale', key)
}

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped lang="less">
.admin-header {
  background: var(--header-bg, #fff);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  border-bottom: 1px solid var(--border-color, #f0f0f0);
  position: sticky;
  top: 0;
  z-index: 9;
  transition: background-color 0.3s;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-breadcrumb {
  line-height: 64px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.lang-btn {
  .lang-text {
    margin-left: 4px;
    font-size: 12px;
  }
}

.user-info {
  color: var(--text-color, #333);
}
</style>
