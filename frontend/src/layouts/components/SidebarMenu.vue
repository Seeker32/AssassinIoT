<template>
  <div class="sidebar">
    <!-- Logo -->
    <div class="sidebar-logo">
      <el-icon :size="24" color="#409eff">
        <Cpu />
      </el-icon>
      <span v-show="!appStore.sidebarCollapsed" class="logo-text">AssassinIoT</span>
    </div>

    <!-- Menu -->
    <el-menu
      :default-active="route.path"
      :collapse="appStore.sidebarCollapsed"
      :router="true"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409eff"
    >
      <template v-for="item in menuRoutes" :key="item.path">
        <el-menu-item :index="resolvePath(item)">
          <el-icon v-if="item.meta?.icon">
            <component :is="item.meta.icon as string" />
          </el-icon>
          <template #title>{{ item.meta?.title }}</template>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { routes } from '@/router/routes'

const route = useRoute()
const appStore = useAppStore()

const menuRoutes = computed(() => {
  const adminRoute = routes.find(r => r.path === '/')
  return adminRoute?.children?.filter(r => !r.meta?.hidden) || []
})

function resolvePath(item: (typeof menuRoutes.value)[number]): string {
  // Build the full path for the menu item
  return '/' + item.path
}
</script>

<style scoped lang="scss">
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
  white-space: nowrap;
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.el-menu {
  border-right: none;
  flex: 1;
}
</style>
