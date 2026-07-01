<template>
  <div class="thing-model-detail">
    <div class="page-header">
      <el-button text @click="router.back()">
        <el-icon><ArrowLeft /></el-icon> 返回
      </el-button>
      <h2>物模型详情</h2>
    </div>

    <el-card v-loading="loading">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ model?.id }}</el-descriptions-item>
        <el-descriptions-item label="模型标识">{{ model?.model_key }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ model?.name }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ model?.category }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ model?.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag
            :type="
              model?.status === 'active'
                ? 'success'
                : model?.status === 'deprecated'
                  ? 'warning'
                  : 'info'
            "
          >
            {{ model?.status === 'active' ? '启用' : model?.status === 'deprecated' ? '已弃用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ model?.description || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <h3>属性定义 (Properties)</h3>
      <el-empty v-if="!propertyList.length" description="暂无属性定义" />
      <el-table v-else :data="propertyList" stripe>
        <el-table-column prop="identifier" label="标识符" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="data_type" label="数据类型" width="100" />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="mode" label="读写模式" width="80" />
        <el-table-column prop="description" label="描述" />
      </el-table>

      <el-divider />

      <h3>服务定义 (Services)</h3>
      <el-empty v-if="!serviceList.length" description="暂无服务定义" />
      <el-table v-else :data="serviceList" stripe>
        <el-table-column prop="identifier" label="标识符" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="description" label="描述" />
      </el-table>

      <el-divider />

      <h3>事件定义 (Events)</h3>
      <el-empty v-if="!eventList.length" description="暂无事件定义" />
      <el-table v-else :data="eventList" stripe>
        <el-table-column prop="identifier" label="标识符" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="data_type" label="数据类型" width="100" />
        <el-table-column prop="description" label="描述" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchThingModel } from '@/api/thing-models'
import type { ThingModel, PropertyDef, ServiceDef, EventDef } from '@/api/types/thing-model'

const route = useRoute()
const router = useRouter()
const model = ref<ThingModel | null>(null)
const loading = ref(false)

const propertyList = computed<PropertyDef[]>(() => {
  if (!model.value) return []
  return Object.values(model.value.properties)
})

const serviceList = computed<ServiceDef[]>(() => {
  if (!model.value) return []
  return Object.values(model.value.services)
})

const eventList = computed<EventDef[]>(() => {
  if (!model.value) return []
  return Object.values(model.value.events)
})

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id || isNaN(id)) return

  loading.value = true
  try {
    model.value = await fetchThingModel(id)
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
</style>
