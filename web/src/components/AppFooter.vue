<template>
  <a-layout-footer class="app-footer">
    <div class="footer-inner">
      <div class="footer-content">
        <div class="footer-copyright">
          {{ appStore.siteConfig.copyright || `© ${new Date().getFullYear()} ${appStore.siteConfig.siteName}` }}
        </div>
        <div class="footer-links">
          <router-link to="/">{{ $t('nav.home') }}</router-link>
          <a-divider type="vertical" />
          <router-link to="/?tab=categories">{{ $t('nav.categories') }}</router-link>
          <a-divider type="vertical" />
          <router-link to="/?tab=tags">{{ $t('nav.tags') }}</router-link>
        </div>
        <div v-if="appStore.siteConfig.icp" class="footer-icp">
          <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener">{{ appStore.siteConfig.icp }}</a>
        </div>
      </div>
    </div>
  </a-layout-footer>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
</script>

<style scoped lang="less">
.app-footer {
  background: var(--footer-bg);
  border-top: 1px solid var(--border-color);
  padding: 32px 16px;
  text-align: center;
  transition: background-color 0.3s;
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.footer-copyright {
  color: var(--text-secondary);
  font-size: 14px;
}

.footer-links {
  a {
    color: var(--text-secondary);
    font-size: 13px;
    text-decoration: none;
    &:hover {
      color: var(--primary-color);
    }
  }
}

.footer-icp {
  a {
    color: var(--text-secondary);
    font-size: 12px;
    text-decoration: none;
    &:hover {
      color: var(--primary-color);
    }
  }
}
</style>
