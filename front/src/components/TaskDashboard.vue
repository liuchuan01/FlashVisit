<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>{{ taskData.name }}</h2>
      <div class="task-controls">
        <el-tag :type="statusType">{{ taskData.status }}</el-tag>
        <el-button 
          type="danger" 
          :disabled="taskData.status !== 'running'"
          @click="stopTask"
        >
          停止任务
        </el-button>
      </div>
    </div>

    <!-- 进度条 -->
    <el-card class="progress-card">
      <template #header>
        <div class="card-header">
          <span>任务进度</span>
          <span class="time-remaining">剩余时间: {{ formatTimeRemaining }}</span>
        </div>
      </template>
      <el-progress 
        :percentage="progressPercentage" 
        :status="progressStatus"
      ></el-progress>
    </el-card>

    <!-- URL卡片网格 -->
    <div class="url-grid">
      <el-card 
        v-for="url in urlStatuses" 
        :key="url.url" 
        class="url-card"
      >
        <template #header>
          <div class="url-header">
            <el-tooltip :content="url.url" placement="top">
              <span class="url-text">{{ truncateUrl(url.url) }}</span>
            </el-tooltip>
          </div>
        </template>
        <div class="url-stats">
          <div class="stat-item">
            <span class="label">访问次数</span>
            <span class="value">{{ url.visitCount }}</span>
          </div>
          <div class="stat-item">
            <span class="label">最后访问</span>
            <span class="value">{{ formatLastVisit(url.lastVisitTime) }}</span>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

export default {
  name: 'TaskDashboard',
  props: {
    taskId: {
      type: String,
      required: true
    }
  },
  emits: ['refresh-tasks'],
  setup(props, { emit }) {
    const taskData = ref({
      name: '',
      status: '',
      createdAt: '',
      timeout: 0
    })
    const urlStatuses = ref([])
    let refreshInterval = null

    const statusType = computed(() => {
      switch (taskData.value.status) {
        case 'running': return 'success'
        case 'stopped': return 'danger'
        case 'completed': return 'info'
        default: return 'warning'
      }
    })

    const progressPercentage = computed(() => {
      if (!taskData.value.createdAt || !taskData.value.timeout) return 0
      const startTime = new Date(taskData.value.createdAt).getTime()
      const currentTime = new Date().getTime()
      const elapsedSeconds = (currentTime - startTime) / 1000
      return Math.min(Math.round((elapsedSeconds / taskData.value.timeout) * 100), 100)
    })

    const progressStatus = computed(() => {
      if (taskData.value.status === 'stopped') return 'exception'
      if (progressPercentage.value >= 100) return 'success'
      return ''
    })

    const formatTimeRemaining = computed(() => {
      if (!taskData.value.createdAt || !taskData.value.timeout) return '未知'
      const startTime = new Date(taskData.value.createdAt).getTime()
      const currentTime = new Date().getTime()
      const elapsedSeconds = (currentTime - startTime) / 1000
      const remainingSeconds = Math.max(taskData.value.timeout - elapsedSeconds, 0)
      
      const hours = Math.floor(remainingSeconds / 3600)
      const minutes = Math.floor((remainingSeconds % 3600) / 60)
      const seconds = Math.floor(remainingSeconds % 60)
      
      return `${hours}小时${minutes}分${seconds}秒`
    })

    const fetchTaskData = async () => {
      try {
        const [taskResponse, urlResponse] = await Promise.all([
          axios.get(`/api/tasks`),
          axios.get(`/api/tasks/${props.taskId}/urls`)
        ])
        
        taskData.value = taskResponse.data.find(task => task.id === props.taskId) || {}
        urlStatuses.value = urlResponse.data
      } catch (error) {
        ElMessage.error('获取任务数据失败：' + error.message)
      }
    }

    const stopTask = async () => {
      try {
        await axios.post(`/api/tasks/${props.taskId}/stop`)
        ElMessage.success('任务已停止')
        emit('refresh-tasks')
      } catch (error) {
        ElMessage.error('停止任务失败：' + error.message)
      }
    }

    const truncateUrl = (url) => {
      if (url.length > 40) {
        return url.substring(0, 37) + '...'
      }
      return url
    }

    const formatLastVisit = (time) => {
      if (!time) return '未访问'
      const date = new Date(time)
      return date.toLocaleString()
    }

    onMounted(() => {
      fetchTaskData()
      refreshInterval = setInterval(fetchTaskData, 2000)
    })

    onUnmounted(() => {
      if (refreshInterval) {
        clearInterval(refreshInterval)
      }
    })

    return {
      taskData,
      urlStatuses,
      statusType,
      progressPercentage,
      progressStatus,
      formatTimeRemaining,
      stopTask,
      truncateUrl,
      formatLastVisit
    }
  }
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.task-controls {
  display: flex;
  gap: 15px;
  align-items: center;
}

.progress-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.time-remaining {
  font-size: 14px;
  color: #666;
}

.url-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.url-card {
  height: 100%;
}

.url-header {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.url-text {
  font-size: 14px;
  color: #409EFF;
  cursor: pointer;
}

.url-stats {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label {
  color: #666;
  font-size: 14px;
}

.value {
  font-weight: bold;
  color: #333;
}
</style>
