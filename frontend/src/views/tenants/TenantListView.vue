<template>
  <div class="tenant-list">
    <h2>租户管理</h2>

    <el-card>
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索租户..."
          clearable
          style="width: 300px"
          @keyup.enter="loadTenants"
        />
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增租户
        </el-button>
      </div>

      <el-table :data="tenants" v-loading="loading" stripe>
        <el-table-column prop="tenant_key" label="租户标识" min-width="140" />
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="router.push(`/tenants/${row.id}`)">
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
          @change="loadTenants"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchTenants, deleteTenant } from '@/api/tenants'
import type { Tenant } from '@/api/types/tenant'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const tenants = ref<Tenant[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

async function loadTenants() {
  loading.value = true
  try {
    const res = await fetchTenants({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
    })
    tenants.value = res.data
    total.value = res.total
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  // TODO: open create dialog or navigate to create page
  router.push('/tenants/new')
}

function handleEdit(row: Tenant) {
  // TODO: open edit dialog
  router.push(`/tenants/${row.id}`)
}

async function handleDelete(row: Tenant) {
  try {
    await ElMessageBox.confirm(
      `确认删除租户「${row.name}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' },
    )
    await deleteTenant(row.id)
    loadTenants()
  } catch {
    // cancelled or error
  }
}

onMounted(loadTenants)
</script>
