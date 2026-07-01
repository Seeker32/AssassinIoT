<template>
  <div class="device-list">
    <h2>设备管理</h2>

    <el-card>
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索设备..."
          clearable
          style="width: 300px"
          @keyup.enter="loadList"
        />
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          注册设备
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="设备ID" min-width="180" show-overflow-tooltip />
        <el-table-column prop="device_name" label="设备名称" min-width="160" />
        <el-table-column prop="tenant_key" label="租户" width="120" />
        <el-table-column prop="model_key" label="物模型" width="140" />
        <el-table-column prop="online" label="在线状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.online ? 'success' : 'info'">
              {{ row.online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="
                row.status === 'active'
                  ? 'success'
                  : row.status === 'inactive'
                    ? 'warning'
                    : 'info'
              "
            >
              {{ row.status === 'active' ? '活跃' : row.status === 'inactive' ? '待激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="firmware_ver" label="固件版本" width="120" />
        <el-table-column prop="last_seen" label="最后在线" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="router.push(`/devices/${row.id}`)">
              查看
            </el-button>
            <el-button text type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button text type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="loadList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchDevices, deleteDevice } from '@/api/devices'
import type { Device } from '@/api/types/device'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const list = ref<Device[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

async function loadList() {
  loading.value = true
  try {
    const res = await fetchDevices({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
    })
    list.value = res.data
    total.value = res.total
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  // TODO: open create dialog
}

function handleEdit(row: Device) {
  router.push(`/devices/${row.id}`)
}

async function handleDelete(row: Device) {
  try {
    await ElMessageBox.confirm(
      `确认删除设备「${row.device_name}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' },
    )
    await deleteDevice(row.id)
    loadList()
  } catch {
    // cancelled or error
  }
}

onMounted(loadList)
</script>
