import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const siteConfig = ref({
    siteName: 'GoBlog',
    siteLogo: '',
    copyright: '',
    icp: '',
    registerEnabled: true,
  })

  function setSiteConfig(config: Record<string, any>) {
    siteConfig.value = { ...siteConfig.value, ...config }
  }

  return { siteConfig, setSiteConfig }
})
