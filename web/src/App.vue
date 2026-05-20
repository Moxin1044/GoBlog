<template>
  <a-config-provider :locale="antdLocale" :theme="themeConfig">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { theme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import enUS from 'ant-design-vue/es/locale/en_US'
import { useThemeStore } from '@/stores/theme'

const { locale } = useI18n()
const themeStore = useThemeStore()

const antdLocale = computed(() => locale.value === 'zh-CN' ? zhCN : enUS)

const themeConfig = computed(() => ({
  algorithm: themeStore.isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 6,
  },
}))
</script>
