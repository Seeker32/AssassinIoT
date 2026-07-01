<template>
  <div class="tenant-detail">
    <div class="page-header">
      <el-button text @click="router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>租户详情</h2>
    </div>

    <el-card v-loading="loading">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ tenant?.id }}</el-descriptions-item>
        <el-descriptions-item label="租户标识">{{ tenant?.tenant_key }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ tenant?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="tenant?.status === 'active' ? 'success' : 'info'">
            {{ tenant?.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ tenant?.description || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ tenant?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ tenant?.updated_at }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchTenant } from '@/api/tenants'
import type { Tenant } from '@/api/types/tenant'

const route = useRoute()
const router = useRouter()
const tenant = ref<Tenant | null>(null)
const loading = ref(false)

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id || isNaN(id)) return

  loading.value = true
  try {
    tenant.value = await fetchTenant(id)
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
})
</script>

<style scoped lang="scss">
.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;

  h2 {
    margin-bottom: 0;
  }
}
</style>
