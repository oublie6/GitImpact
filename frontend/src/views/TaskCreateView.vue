<template>
  <div>
    <h3>新建任务</h3>
    <el-input v-model="task.task_name" placeholder="任务名" />
    <el-input v-model="task.old_ref" placeholder="old_ref" />
    <el-input v-model="task.new_ref" placeholder="new_ref" />
    <el-button @click="submit">提交</el-button>
  </div>
</template>

<script setup lang="ts">
// 任务创建页当前使用固定仓库 ID 作为最小演示路径，尚未做仓库选择器。
import { reactive } from 'vue'
import http from '../api/http'

const task = reactive({
  task_name: '',
  mode: 'same_repo_commits',
  old_repo_id: 1,
  new_repo_id: 1,
  old_ref: 'main~1',
  new_ref: 'main',
  generate_markdown: true,
  generate_structured: true,
  custom_focus: '',
  remark: ''
})

const submit = async () => {
  await http.post('/api/tasks', task)
}
</script>
