<template>
  <a-layout class="admin-layout" :class="{ dark: themeStore.isDark }">
    <!-- Desktop Sidebar -->
    <AdminSidebar
      v-if="!isMobile"
      :collapsed="collapsed"
      @toggle="collapsed = !collapsed"
    />
    <!-- Mobile Drawer Sidebar -->
    <a-drawer
      v-if="isMobile"
      :open="!collapsed"
      placement="left"
      :closable="false"
      :width="220"
      :body-style="{ padding: 0 }"
      @close="collapsed = true"
    >
      <AdminSidebar :collapsed="false" @toggle="collapsed = !collapsed" />
    </a-drawer>
    <a-layout :style="{ marginLeft: isMobile ? '0' : (collapsed ? '80px' : '220px'), transition: 'margin-left 0.2s' }">
      <AdminHeader :collapsed="collapsed" @toggle="collapsed = !collapsed" />
      <a-layout-content class="admin-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import AdminSidebar from '@/components/AdminSidebar.vue'
import AdminHeader from '@/components/AdminHeader.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const collapsed = ref(false)
const windowWidth = ref(window.innerWidth)

const isMobile = computed(() => windowWidth.value < 768)

function handleResize() {
  windowWidth.value = window.innerWidth
  if (isMobile.value) {
    collapsed.value = true
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  handleResize()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="less">
.admin-layout {
  min-height: 100vh;
}

.admin-content {
  margin: 24px;
  padding: 24px;
  border-radius: 8px;
  min-height: 280px;
}

@media (max-width: 768px) {
  .admin-content {
    margin: 12px;
    padding: 12px;
  }
}
</style>
