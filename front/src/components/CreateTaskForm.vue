<template>
  <el-form :model="form" label-width="120px" @submit.prevent="handleSubmit" class="create-task-form">
    <el-form-item label="任务名称" required>
      <el-input
        v-model="form.name"
        placeholder="请输入任务名称"
        class="custom-input"
      ></el-input>
    </el-form-item>
    
    <el-form-item label="URL列表" required>
      <el-input
        type="textarea"
        v-model="form.urls"
        :rows="6"
        placeholder="请输入URL列表，每行一个"
        class="custom-textarea"
      ></el-input>
    </el-form-item>
    
    <el-form-item label="时间设置" class="time-settings">
      <div class="time-inputs">
        <div class="time-input-group">
          <label>最小间隔(秒)</label>
          <el-input-number
            v-model="form.minInterval"
            :min="1"
            :max="form.maxInterval"
            class="custom-number-input"
          ></el-input-number>
        </div>
        <div class="time-input-group">
          <label>最大间隔(秒)</label>
          <el-input-number
            v-model="form.maxInterval"
            :min="form.minInterval"
            class="custom-number-input"
          ></el-input-number>
        </div>
        <div class="time-input-group">
          <label>超时时间(秒)</label>
          <el-input-number
            v-model="form.timeout"
            :min="1"
            class="custom-number-input"
          ></el-input-number>
        </div>
      </div>
    </el-form-item>

    <el-form-item class="form-footer">
      <el-button type="primary" @click="handleSubmit" class="submit-button">创建任务</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

export default {
  name: 'CreateTaskForm',
  emits: ['task-created'],
  setup(props, { emit }) {
    const form = ref({
      name: '',
      urls: '',
      minInterval: 65,
      maxInterval: 95,
      timeout: 86400
    })

    const handleSubmit = async () => {
      try {
        if (!form.value.name || !form.value.urls) {
          ElMessage.warning('请填写任务名称和URL列表')
          return
        }

        const urls = form.value.urls.split('\n').filter(url => url.trim())
        if (urls.length === 0) {
          ElMessage.warning('请至少输入一个有效的URL')
          return
        }

        await axios.post('/api/tasks', {
          name: form.value.name,
          urls: urls,
          minInterval: form.value.minInterval,
          maxInterval: form.value.maxInterval,
          timeout: form.value.timeout
        })

        ElMessage.success('任务创建成功')
        emit('task-created')
        
        // 重置表单
        form.value = {
          name: '',
          urls: '',
          minInterval: 65,
          maxInterval: 95,
          timeout: 86400
        }
      } catch (error) {
        ElMessage.error('创建任务失败：' + error.message)
      }
    }

    return {
      form,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.create-task-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.custom-input {
  width: 100%;
  --el-input-border-radius: 8px;
}

.custom-input :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  border: none;
  box-shadow: none;
  transition: all 0.3s ease;
}

.custom-input :deep(.el-input__wrapper:hover),
.custom-input :deep(.el-input__wrapper.is-focus) {
  background-color: #eef0f6;
  box-shadow: 0 0 0 1px var(--el-color-primary) inset;
}

.custom-textarea {
  width: 100%;
  --el-input-border-radius: 8px;
}

.custom-textarea :deep(.el-textarea__inner) {
  background-color: #f5f7fa;
  border: none;
  box-shadow: none;
  transition: all 0.3s ease;
  min-height: 150px !important;
  padding: 12px;
}

.custom-textarea :deep(.el-textarea__inner:hover),
.custom-textarea :deep(.el-textarea__inner:focus) {
  background-color: #eef0f6;
  box-shadow: 0 0 0 1px var(--el-color-primary) inset;
}

.time-settings {
  margin-top: 24px;
}

.time-inputs {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-top: 8px;
}

.time-input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.time-input-group label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.custom-number-input {
  width: 100%;
}

.custom-number-input :deep(.el-input-number__decrease),
.custom-number-input :deep(.el-input-number__increase) {
  background-color: transparent;
  border: none;
  color: var(--el-color-primary);
}

.custom-number-input :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  border: none;
  box-shadow: none;
  border-radius: 8px;
}

.form-footer {
  margin-top: 32px;
  text-align: right;
}

.submit-button {
  width: 120px;
  height: 40px;
  border-radius: 20px;
  font-weight: 500;
  background: linear-gradient(to right, #409EFF, #36a3ff);
  border: none;
  transition: all 0.3s ease;
}

.submit-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #333;
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

:deep(.el-form-item) {
  margin-bottom: 24px;
}
</style>
