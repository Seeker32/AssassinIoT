<template>
  <div class="device-detail">
    <div class="page-header">
      <el-button text @click="router.back()">
        <el-icon><ArrowLeft /></el-icon> 返回
      </el-button>
      <h2>设备详情</h2>
    </div>

    <el-card v-loading="loading">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="设备ID">{{ device?.id }}</el-descriptions-item>
        <el-descriptions-item label="设备名称">{{ device?.device_name }}</el-descriptions-item>
        <el-descriptions-item label="租户">{{ device?.tenant_key }}</el-descriptions-item>
        <el-descriptions-item label="物模型">{{ device?.model_key }}</el-descriptions-item>
        <el-descriptions-item label="固件版本">{{ device?.firmware_ver || '-' }}</el-descriptions-item>
        <el-descriptions-item label="在线状态">
          <el-tag :type="device?.online ? 'success' : 'info'">
            {{ device?.online ? '在线' : '离线' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="设备状态">
          <el-tag
            :type="
              device?.status === 'active'
                ? 'success'
                : device?.status === 'inactive'
                  ? 'warning'
                  : 'info'
            "
          >
            {{ device?.status === 'active' ? '活跃' : device?.status === 'inactive' ? '待激活' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="最后在线">{{ device?.last_seen || '-' }}</el-descriptions-item>
        <el-descriptions-item label="接入密钥" :span="2">
          <code>{{ device?.access_key }}</code>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ device?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ device?.updated_at }}</el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <h3>遥测数据</h3>
      <el-empty description="遥测数据功能开发中" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchDevice } from '@/api/devices'
import type { Device } from '@/api/types/device'

const route = useRoute()
const router = useRouter()
const device = ref<Device | null>(null)
const loading = ref(false)

onMounted(async () => {
  const id = route.params.id as string
  if (!id) return

  loading.value = true
  try {
    device.value = await fetchDevice(id)
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
  h2 { margin-bottom: 0; }
}

h3 {
  font-size: 16px;
  margin-bottom: 12px;
}

code {
  background: var(--el-fill-color-light);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 13px;
}
</style>
