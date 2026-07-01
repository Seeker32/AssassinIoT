<template>
  <el-container class="admin-layout">
    <!-- Sidebar -->
    <el-aside :width="appStore.sidebarCollapsed ? '64px' : '220px'" class="admin-aside">
      <SidebarMenu />
    </el-aside>

    <!-- Main area -->
    <el-container>
      <!-- Header -->
      <el-header class="admin-header" height="56px">
        <HeaderBar />
      </el-header>

      <!-- Content -->
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import SidebarMenu from './components/SidebarMenu.vue'
import HeaderBar from './components/HeaderBar.vue'

const appStore = useAppStore()
</script>

<style scoped lang="scss">
.admin-layout {
  height: 100vh;
}

.admin-aside {
  transition: width 0.3s;
  background-color: #304156;
  overflow: hidden;
}

.admin-header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid var(--el-border-color-light);
  padding: 0 20px;
  height: 56px;
}

.admin-main {
  background-color: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}
</style>
