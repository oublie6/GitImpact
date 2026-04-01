<template>
  <div>
    <h3>任务详情</h3>
    <el-button @click="load">加载</el-button>
    <pre>{{ detail }}</pre>
    <pre>{{ report }}</pre>
  </div>
</template>

<script setup lang="ts">
// 详情页顺序加载任务元数据和报告内容，方便排查任务执行结果。
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import http from '../api/http'

const route = useRoute()
const detail = ref({})
const report = ref({})

const load = async () => {
  detail.value = (await http.get(`/api/tasks/${route.params.id}`)).data.data
  report.value = (await http.get(`/api/tasks/${route.params.id}/report`)).data.data
}
</script>
