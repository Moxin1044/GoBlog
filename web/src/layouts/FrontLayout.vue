<template>
  <div class="front-layout" :class="{ dark: themeStore.isDark }">
    <AppHeader />
    <div class="front-content-wrapper">
      <div class="front-main">
        <router-view />
      </div>
      <AppSidebar class="front-sidebar" />
    </div>
    <AppFooter />
    <ChatPanel />
    <!-- Back to Top -->
    <a-back-top :style="{ right: '24px', bottom: '96px' }" />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import AppSidebar from '@/components/AppSidebar.vue'
import ChatPanel from '@/components/ChatPanel.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
</script>

<style scoped lang="less">
.front-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-color);
  color: var(--text-color);
  transition: background-color 0.3s, color 0.3s;
}

.front-content-wrapper {
  flex: 1;
  display: flex;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 24px 16px;
  gap: 24px;
}

.front-main {
  flex: 1;
  min-width: 0;
}

.front-sidebar {
  width: 300px;
  flex-shrink: 0;
}

@media (max-width: 992px) {
  .front-sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .front-content-wrapper {
    padding: 16px 12px;
  }
}
</style>
