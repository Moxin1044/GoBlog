import { useThemeStore } from '@/stores/theme'
import { computed } from 'vue'

export function useTheme() {
  const themeStore = useThemeStore()
  const isDark = computed(() => themeStore.isDark)
  const toggleTheme = () => themeStore.toggleTheme()
  return { isDark, toggleTheme }
}
