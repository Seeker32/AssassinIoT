import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const currentTenantKey = ref<string>('')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTenant(tenantKey: string) {
    currentTenantKey.value = tenantKey
  }

  return {
    sidebarCollapsed,
    currentTenantKey,
    toggleSidebar,
    setTenant,
  }
})
