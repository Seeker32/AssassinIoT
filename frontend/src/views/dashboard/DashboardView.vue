<template>
  <div class="dashboard">
    <h2>仪表盘</h2>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>设备总数</span>
              <el-icon :size="20" color="#409eff"><Monitor /></el-icon>
            </div>
          </template>
          <div class="stat-value">--</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>在线设备</span>
              <el-icon :size="20" color="#67c23a"><CircleCheck /></el-icon>
            </div>
          </template>
          <div class="stat-value stat-online">--</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>租户数</span>
              <el-icon :size="20" color="#e6a23c"><OfficeBuilding /></el-icon>
            </div>
          </template>
          <div class="stat-value">--</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>物模型数</span>
              <el-icon :size="20" color="#909399"><Setting /></el-icon>
            </div>
          </template>
          <div class="stat-value">--</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped lang="scss">
.dashboard {
  h2 {
    margin-bottom: 20px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  text-align: center;
  padding: 10px 0;
}

.stat-online {
  color: #67c23a;
}
</style>
