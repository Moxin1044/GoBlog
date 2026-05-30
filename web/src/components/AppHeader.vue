<template>
  <a-layout-header class="app-header">
    <div class="header-inner">
      <div class="header-left">
        <router-link to="/" class="logo-link">
          <img v-if="appStore.siteConfig.siteLogo" :src="appStore.siteConfig.siteLogo" alt="Logo" class="logo-img" />
          <span class="site-name">{{ appStore.siteConfig.siteName }}</span>
        </router-link>
        <a-menu
          mode="horizontal"
          :selected-keys="selectedKeys"
          class="nav-menu"
          :overflowedIndicator="null"
          @click="handleMenuClick"
        >
          <a-menu-item key="home">{{ $t('nav.home') }}</a-menu-item>
          <template v-for="nav in flatNavigations" :key="nav.id">
            <a-menu-item :key="nav.id" @click="handleNavClick(nav)">
              {{ getNavName(nav) }}
            </a-menu-item>
          </template>
        </a-menu>
      </div>
      <div class="header-right">
        <a-input-search
          v-model:value="searchText"
          :placeholder="$t('common.search')"
          class="header-search"
          @search="handleSearch"
          allow-clear
        />
        <a-dropdown>
          <a-button type="text" size="small" class="lang-btn">
            {{ currentLangLabel }}
          </a-button>
          <template #overlay>
            <a-menu @click="handleLangChange">
              <a-menu-item key="zh-CN">中文</a-menu-item>
              <a-menu-item key="en-US">English</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
        <ThemeToggle />
        <template v-if="userStore.isLoggedIn()">
          <a-dropdown>
            <div class="user-avatar-wrapper">
              <a-avatar :style="{ cursor: 'pointer' }" :src="userAvatar">
                {{ userStore.username.charAt(0).toUpperCase() }}
              </a-avatar>
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile" @click="$router.push('/profile')">
                  <UserOutlined /> {{ $t('nav.profile') }}
                </a-menu-item>
                <a-menu-item v-if="userStore.isAdmin()" key="admin" @click="$router.push('/admin')">
                  <SettingOutlined /> {{ $t('nav.admin') }}
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined /> {{ $t('nav.logout') }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </template>
        <template v-else>
          <a-button type="primary" @click="$router.push('/login')">{{ $t('nav.login') }}</a-button>
        </template>
        <!-- Mobile Menu Button -->
        <a-button type="text" class="mobile-menu-btn show-mobile-only" @click="mobileMenuVisible = true">
          <MenuOutlined />
        </a-button>
      </div>
    </div>

    <!-- Mobile Drawer -->
    <a-drawer
      v-model:open="mobileMenuVisible"
      placement="left"
      :title="appStore.siteConfig.siteName"
      :closable="true"
      :width="280"
    >
      <a-menu mode="vertical" :selected-keys="selectedKeys" @click="handleMobileMenuClick">
        <a-menu-item key="home" @click="router.push('/')">
          <HomeOutlined /> {{ $t('nav.home') }}
        </a-menu-item>
        <template v-for="nav in flatNavigations" :key="nav.id">
          <a-menu-item :key="nav.id" @click="handleNavClick(nav)">
            <AppstoreOutlined /> {{ getNavName(nav) }}
          </a-menu-item>
        </template>
      </a-menu>
      <div class="mobile-search">
        <a-input-search
          v-model:value="searchText"
          :placeholder="$t('common.search')"
          @search="handleSearch"
          allow-clear
        />
      </div>
    </a-drawer>
  </a-layout-header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  UserOutlined, SettingOutlined, LogoutOutlined,
  MenuOutlined, HomeOutlined, AppstoreOutlined, TagsOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import ThemeToggle from './ThemeToggle.vue'
import { getNavigations } from '@/api/article'

const router = useRouter()
const route = useRoute()
const { locale } = useI18n()
const userStore = useUserStore()
const appStore = useAppStore()

const searchText = ref('')
const mobileMenuVisible = ref(false)
const navigations = ref<any[]>([])

const flatNavigations = computed(() => {
  const result: any[] = []
  for (const nav of navigations.value) {
    result.push(nav)
    if (nav.children && nav.children.length > 0) {
      for (const child of nav.children) {
        result.push(child)
      }
    }
  }
  return result
})

const selectedKeys = computed(() => {
  if (route.path === '/') return ['home']
  return []
})

const currentLangLabel = computed(() => locale.value === 'zh-CN' ? '中文' : 'EN')

const userAvatar = computed(() => {
  return '' // Could be extended to read from user store
})

const getNavName = (nav: any) => {
  if (locale.value === 'en-US' && nav.name_en) {
    return nav.name_en
  }
  return nav.name
}

const fetchNavigations = async () => {
  try {
    const res = await getNavigations()
    navigations.value = res.data || []
  } catch {
    // Ignore error
  }
}

function handleNavClick(nav: any) {
  if (nav.type === 'category' && nav.category_id) {
    router.push(`/category/${nav.category_id}`)
  } else if (nav.type === 'link' && nav.link) {
    if (nav.new_tab) {
      window.open(nav.link, '_blank')
    } else {
      window.location.href = nav.link
    }
  } else {
    // Custom type, treat as custom link or internal link
    if (nav.link) {
      if (nav.new_tab) {
        window.open(nav.link, '_blank')
      } else {
        if (nav.link.startsWith('http')) {
          window.location.href = nav.link
        } else {
          router.push(nav.link)
        }
      }
    }
  }
  mobileMenuVisible.value = false
}

function handleMenuClick({ key }: { key: string }) {
  if (key === 'home') router.push('/')
}

function handleMobileMenuClick({ key }: { key: string }) {
  handleMenuClick({ key })
  mobileMenuVisible.value = false
}

function handleSearch(value: string) {
  if (value.trim()) {
    router.push(`/?search=${encodeURIComponent(value.trim())}`)
    mobileMenuVisible.value = false
  }
}

function handleLangChange({ key }: { key: string }) {
  locale.value = key
  localStorage.setItem('locale', key)
}

function handleLogout() {
  userStore.logout()
  router.push('/')
}

onMounted(() => {
  fetchNavigations()
})
</script>

<style scoped lang="less">
.app-header {
  background: var(--header-bg);
  border-bottom: 1px solid var(--border-color);
  padding: 0 24px;
  height: 64px;
  line-height: 64px;
  position: sticky;
  top: 0;
  z-index: 100;
  transition: background-color 0.3s;
}

.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: var(--text-color);
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
}

.logo-img {
  height: 32px;
}

.nav-menu {
  background: transparent;
  border-bottom: none;
  line-height: 62px;
  min-width: 0;

  :deep(.ant-menu-overflow) {
    display: flex;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-search {
  width: 200px;
}

.lang-btn {
  font-size: 13px;
}

.user-avatar-wrapper {
  display: flex;
  align-items: center;
}

.mobile-menu-btn {
  font-size: 18px;
  display: none;
}

.mobile-search {
  margin-top: 16px;
  padding: 0 16px;
}

@media (max-width: 768px) {
  .nav-menu {
    display: none;
  }
  .header-search {
    width: 120px;
  }
  .mobile-menu-btn {
    display: inline-flex;
  }
}

@media (max-width: 576px) {
  .header-search {
    display: none;
  }
  .site-name {
    font-size: 16px;
  }
}
</style>
