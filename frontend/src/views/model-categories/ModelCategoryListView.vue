<template>
  <div class="model-category-list">
    <h2>物模型分类</h2>

    <el-card>
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索分类..."
          clearable
          style="width: 300px"
          @keyup.enter="loadList"
        />
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增分类
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="category_key" label="分类标识" min-width="140" />
        <el-table-column prop="display_name" label="展示名称" min-width="160" />
        <el-table-column prop="icon" label="图标" width="80" />
        <el-table-column prop="sort_order" label="排序" width="80" />
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
import { fetchModelCategories, deleteModelCategory } from '@/api/model-categories'
import type { ModelCategory } from '@/api/types/model-category'
import { ElMessageBox } from 'element-plus'

const list = ref<ModelCategory[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

async function loadList() {
  loading.value = true
  try {
    const res = await fetchModelCategories({
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

function handleEdit(_row: ModelCategory) {
  // TODO: open edit dialog
}

async function handleDelete(row: ModelCategory) {
  try {
    await ElMessageBox.confirm(
      `确认删除分类「${row.display_name}」吗？`,
      '删除确认',
      { type: 'warning' },
    )
    await deleteModelCategory(row.id)
    loadList()
  } catch {
    // cancelled or error
  }
}

onMounted(loadList)
</script>
