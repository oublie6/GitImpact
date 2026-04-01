<template>
  <div>
    <h3>仓库管理</h3>
    <el-button @click="load">刷新</el-button>
    <pre>{{ repos }}</pre>
  </div>
</template>

<script setup lang="ts">
// 当前仓库页仍是轻量原型，直接展示后端返回的原始列表。
import { onMounted, ref } from 'vue'
import http from '../api/http'

const repos = ref([] as any[])

const load = async () => {
  repos.value = (await http.get('/api/repositories')).data.data
}

onMounted(load)
</script>
