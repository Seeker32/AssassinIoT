<template>
  <div class="thing-model-list">
    <h2>物模型管理</h2>

    <el-card>
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索物模型..."
          clearable
          style="width: 300px"
          @keyup.enter="loadList"
        />
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增物模型
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="model_key" label="模型标识" min-width="140" />
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="category" label="分类" width="140" />
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="
                row.status === 'active'
                  ? 'success'
                  : row.status === 'deprecated'
                    ? 'warning'
                    : 'info'
              "
            >
              {{ row.status === 'active' ? '启用' : row.status === 'deprecated' ? '已弃用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="router.push(`/thing-models/${row.id}`)">
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
import { fetchThingModels, deleteThingModel } from '@/api/thing-models'
import type { ThingModel } from '@/api/types/thing-model'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const list = ref<ThingModel[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

async function loadList() {
  loading.value = true
  try {
    const res = await fetchThingModels({
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

function handleEdit(row: ThingModel) {
  router.push(`/thing-models/${row.id}`)
}

async function handleDelete(row: ThingModel) {
  try {
    await ElMessageBox.confirm(
      `确认删除物模型「${row.name}」吗？`,
      '删除确认',
      { type: 'warning' },
    )
    await deleteThingModel(row.id)
    loadList()
  } catch {
    // cancelled or error
  }
}

onMounted(loadList)
</script>
