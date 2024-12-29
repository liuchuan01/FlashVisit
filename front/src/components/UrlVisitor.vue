<template>
  <div class="url-visitor">
    <el-row :gutter="24">
      <!-- Left Column - Task Creation -->
      <el-col :span="10">
        <el-card class="task-creation" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Plus /></el-icon>
              <span>创建新任务</span>
            </div>
          </template>
          <el-form :model="newTask" label-position="top">
            <el-form-item label="任务名称">
              <el-input 
                v-model="newTask.name" 
                placeholder="请输入任务名称"
                :prefix-icon="Document"
              ></el-input>
            </el-form-item>
            <el-form-item label="URL列表">
              <el-input
                type="textarea"
                v-model="newTask.urls"
                :rows="8"
                placeholder="请输入URL列表，每行一个"
                resize="none"
              ></el-input>
            </el-form-item>
            <el-form-item label="时间设置">
              <el-row :gutter="12">
                <el-col :span="12">
                  <el-form-item label="最小间隔(秒)">
                    <el-input-number
                      v-model="newTask.minInterval"
                      :min="1"
                      controls-position="right"
                      style="width: 100%"
                    ></el-input-number>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="最大间隔(秒)">
                    <el-input-number
                      v-model="newTask.maxInterval"
                      :min="1"
                      controls-position="right"
                      style="width: 100%"
                    ></el-input-number>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="超时时间(秒)">
                <el-input-number
                  v-model="newTask.timeout"
                  :min="1"
                  controls-position="right"
                  style="width: 100%"
                ></el-input-number>
              </el-form-item>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="createTask" style="width: 100%">
                <el-icon><Plus /></el-icon>
                创建任务
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- Right Column - Task List and Stats -->
      <el-col :span="14">
        <!-- Stats Cards -->
        <el-row :gutter="20" class="stats-row">
          <el-col :span="8">
            <el-card shadow="hover" class="stats-card">
              <template #header>
                <div class="stats-header">
                  <el-icon><Timer /></el-icon>
                  <span>运行中任务</span>
                </div>
              </template>
              <div class="stats-value">{{ runningTasks }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover" class="stats-card">
              <template #header>
                <div class="stats-header">
                  <el-icon><Finished /></el-icon>
                  <span>总任务数</span>
                </div>
              </template>
              <div class="stats-value">{{ totalTasks }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover" class="stats-card">
              <template #header>
                <div class="stats-header">
                  <el-icon><Link /></el-icon>
                  <span>总URL数</span>
                </div>
              </template>
              <div class="stats-value">{{ totalUrls }}</div>
            </el-card>
          </el-col>
        </el-row>

        <!-- Task List -->
        <el-card class="task-list" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="header-left">
                <el-icon><List /></el-icon>
                <span>任务列表</span>
              </div>
              <el-button type="primary" @click="refreshTasks" :icon="Refresh">刷新</el-button>
            </div>
          </template>
          <el-table 
            :data="tasks" 
            style="width: 100%"
            :row-class-name="tableRowClassName"
          >
            <el-table-column prop="name" label="任务名称" min-width="120">
              <template #default="scope">
                <div class="task-name">
                  <el-icon><Document /></el-icon>
                  <span>{{ scope.row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="180">
              <template #default="scope">
                <div class="time-cell">
                  <el-icon><Timer /></el-icon>
                  <span>{{ formatDate(scope.row.createdAt) }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button-group>
                  <el-button type="primary" @click="viewTaskDetails(scope.row)" :icon="View">
                    详情
                  </el-button>
                  <el-button type="danger" @click="stopTask(scope.row)" :icon="VideoPlay">
                    停止
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- Task Details Dialog -->
    <el-dialog
      v-model="taskDetailsVisible"
      title="任务详情"
      width="80%"
      destroy-on-close
    >
      <el-tabs>
        <el-tab-pane label="URL列表">
          <el-table 
            :data="currentTaskUrls" 
            style="width: 100%"
            height="400px"
          >
            <el-table-column prop="url" label="URL" min-width="400">
              <template #default="scope">
                <div class="url-cell">
                  <el-icon><Link /></el-icon>
                  <el-link :href="scope.row.url" target="_blank" type="primary">
                    {{ scope.row.url }}
                  </el-link>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="visitCount" label="访问次数" width="120">
              <template #default="scope">
                <el-tag type="success">{{ scope.row.visitCount }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastVisitTime" label="最后访问时间" width="180">
              <template #default="scope">
                <div class="time-cell">
                  <el-icon><Timer /></el-icon>
                  <span>{{ formatDate(scope.row.lastVisitTime) }}</span>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="访问统计">
          <div class="chart-container">
            <v-chart class="chart" :option="visitChartOption" autoresize></v-chart>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'
import { ref, computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import {
  Document,
  Plus,
  Timer,
  Finished,
  Link,
  List,
  Refresh,
  View,
  VideoPlay
} from '@element-plus/icons-vue'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

export default {
  name: 'UrlVisitor',
  components: {
    VChart,
    Document,
    Plus,
    Timer,
    Finished,
    Link,
    List,
    Refresh,
    View,
    VideoPlay
  },
  setup() {
    const newTask = ref({
      name: '',
      urls: '',
      minInterval: 65,
      maxInterval: 95,
      timeout: 86400
    })

    const tasks = ref([])
    const taskDetailsVisible = ref(false)
    const currentTaskUrls = ref([])
    let pollingInterval = null

    // Computed stats
    const runningTasks = computed(() => 
      tasks.value.filter(t => t.status === 'running').length
    )
    const totalTasks = computed(() => tasks.value.length)
    const totalUrls = computed(() => 
      currentTaskUrls.value.reduce((sum, url) => sum + url.visitCount, 0)
    )

    // Chart options
    const visitChartOption = computed(() => ({
      title: {
        text: 'URL访问统计'
      },
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['访问次数']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: currentTaskUrls.value.map(url => url.url.substring(0, 30) + '...')
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '访问次数',
          type: 'line',
          data: currentTaskUrls.value.map(url => url.visitCount),
          smooth: true,
          showSymbol: false,
          areaStyle: {
            opacity: 0.1
          }
        }
      ]
    }))

    return {
      newTask,
      tasks,
      taskDetailsVisible,
      currentTaskUrls,
      runningTasks,
      totalTasks,
      totalUrls,
      visitChartOption,
      Document,
      Plus,
      Timer,
      Finished,
      Link,
      List,
      Refresh,
      View,
      VideoPlay
    }
  },
  methods: {
    async createTask() {
      try {
        await axios.post('/api/tasks', {
          name: this.newTask.name,
          urls: this.newTask.urls.split('\n').filter(url => url.trim()),
          minInterval: this.newTask.minInterval,
          maxInterval: this.newTask.maxInterval,
          timeout: this.newTask.timeout
        })
        this.$message.success('任务创建成功')
        this.refreshTasks()
        // Reset form
        this.newTask = {
          name: '',
          urls: '',
          minInterval: 65,
          maxInterval: 95,
          timeout: 86400
        }
      } catch (error) {
        this.$message.error('任务创建失败：' + error.message)
      }
    },
    async refreshTasks() {
      try {
        const response = await axios.get('/api/tasks')
        this.tasks = response.data
      } catch (error) {
        this.$message.error('获取任务列表失败：' + error.message)
      }
    },
    async viewTaskDetails(task) {
      try {
        const response = await axios.get(`/api/tasks/${task.id}/urls`)
        this.currentTaskUrls = response.data
        this.taskDetailsVisible = true
      } catch (error) {
        this.$message.error('获取任务详情失败：' + error.message)
      }
    },
    async stopTask(task) {
      try {
        await axios.post(`/api/tasks/${task.id}/stop`)
        this.$message.success('任务已停止')
        this.refreshTasks()
      } catch (error) {
        this.$message.error('停止任务失败：' + error.message)
      }
    },
    formatDate(date) {
      return new Date(date).toLocaleString()
    },
    getStatusType(status) {
      const types = {
        running: 'success',
        stopped: 'danger',
        completed: 'info'
      }
      return types[status] || 'info'
    },
    tableRowClassName({ row }) {
      if (row.status === 'running') {
        return 'running-row'
      }
      return ''
    }
  },
  mounted() {
    this.refreshTasks()
    this.pollingInterval = setInterval(this.refreshTasks, 5000)
  },
  beforeUnmount() {
    if (this.pollingInterval) {
      clearInterval(this.pollingInterval)
    }
  }
}
</script>

<style scoped>
.url-visitor {
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
}

.stats-row {
  margin-bottom: 24px;
}

.stats-card {
  height: 140px;
  padding: 8px;
}

.stats-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  color: #606266;
  padding: 12px;
}

.stats-value {
  font-size: 42px;
  font-weight: bold;
  color: var(--el-color-primary);
  text-align: center;
  margin-top: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.task-creation {
  margin-bottom: 24px;
}

:deep(.el-card__header) {
  padding: 16px;
  border-bottom: 2px solid #f0f2f5;
}

:deep(.el-card__body) {
  padding: 24px;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__wrapper) {
  padding: 8px 16px;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-button) {
  padding: 12px 24px;
  font-size: 14px;
}

:deep(.el-table) {
  margin: 16px 0;
}

:deep(.el-table th) {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #606266;
}

:deep(.el-dialog) {
  margin: 0 auto;
  max-width: 1200px;
  width: 90%;
}

.chart-container {
  height: 450px;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
}

:deep(.running-row) {
  background-color: var(--el-color-success-light-9);
}

:deep(.el-card) {
  border-radius: 8px;
  transition: all 0.3s;
}

:deep(.el-card:hover) {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
}

:deep(.el-tag) {
  border-radius: 4px;
}

:deep(.el-tabs__nav-wrap::after) {
  height: 1px;
}
</style>
