<template>
  <div class="app-container">
    <el-container>
      <el-aside width="250px">
        <div class="sidebar">
          <div class="sidebar-header">
            <h2>URL自动访问工具</h2>
            <el-button type="primary" @click="showCreateDialog = true">
              创建任务
            </el-button>
          </div>
          <el-menu
            :default-active="activeTaskId"
            class="task-menu"
            @select="handleTaskSelect"
          >
            <el-menu-item v-for="task in tasks" :key="task.id" :index="task.id">
              <el-icon><Document /></el-icon>
              <span>{{ task.name }}</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-aside>
      <el-main>
        <task-dashboard 
          v-if="activeTaskId" 
          :taskId="activeTaskId"
          @refresh-tasks="fetchTasks"
        />
        <div v-else class="welcome-message">
          <el-empty description="请选择或创建一个任务" />
        </div>
      </el-main>
    </el-container>

    <!-- 创建任务对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建新任务"
      width="50%"
    >
      <create-task-form @task-created="handleTaskCreated" />
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import TaskDashboard from './components/TaskDashboard.vue'
import CreateTaskForm from './components/CreateTaskForm.vue'
import { Document } from '@element-plus/icons-vue'

export default {
  name: 'App',
  components: {
    TaskDashboard,
    CreateTaskForm,
    Document
  },
  setup() {
    const tasks = ref([])
    const activeTaskId = ref('')
    const showCreateDialog = ref(false)

    const fetchTasks = async () => {
      try {
        const response = await axios.get('/api/tasks')
        tasks.value = response.data
      } catch (error) {
        ElMessage.error('获取任务列表失败：' + error.message)
      }
    }

    const handleTaskSelect = (id) => {
      activeTaskId.value = id
    }

    const handleTaskCreated = () => {
      showCreateDialog.value = false
      fetchTasks()
    }

    onMounted(() => {
      fetchTasks()
      // 设置定时刷新
      setInterval(fetchTasks, 5000)
    })

    return {
      tasks,
      activeTaskId,
      showCreateDialog,
      handleTaskSelect,
      handleTaskCreated,
      fetchTasks
    }
  }
}
</script>

<style>
.app-container {
  height: 100vh;
  width: 100vw;
}

.el-container {
  height: 100%;
}

:deep(.el-aside) {
  background-color: #f8f9fa;
  border-right: 1px solid #e6e6e6;
}

.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #dcdfe6;
}

.sidebar-header {
  padding: 24px;
  border-bottom: 1px solid #e6e6e6;
  background-color: #fff;
}

.sidebar-header h2 {
  margin: 0 0 16px 0;
  font-size: 20px;
  font-weight: 500;
  color: #333;
}

.task-menu {
  flex: 1;
  border-right: none;
  background-color: #f8f9fa;
}

.task-menu :deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
  margin: 4px 0;
  border-radius: 0 25px 25px 0;
  margin-right: 16px;
}

.task-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(to right, #ecf5ff 0%, #ecf5ff 90%, transparent 100%);
  color: var(--el-color-primary);
  font-weight: 500;
}

.task-menu :deep(.el-menu-item:hover) {
  background-color: #f0f2f5;
}

.welcome-message {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  background-color: #fff;
  border-radius: 16px;
  margin: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

:deep(.el-dialog__header) {
  margin: 0;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

:deep(.el-dialog__title) {
  font-weight: 500;
  font-size: 18px;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__close) {
  color: #666;
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-button--primary) {
  background: linear-gradient(to right, #409EFF, #36a3ff);
  border: none;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.2);
}

:deep(.el-button--primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
}
</style>
