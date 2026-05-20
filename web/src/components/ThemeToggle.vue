<template>
  <a-button type="text" class="theme-toggle" @click="toggleTheme">
    <template #icon>
      <BulbFilled v-if="themeStore.isDark" />
      <BulbOutlined v-else />
    </template>
  </a-button>
</template>

<script setup lang="ts">
import { BulbOutlined, BulbFilled } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

function toggleTheme() {
  themeStore.toggleTheme()
}
</script>

<style scoped lang="less">
.theme-toggle {
  font-size: 16px;
}
</style>
