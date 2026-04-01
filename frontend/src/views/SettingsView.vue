<template>
  <div>
    <h3>系统设置</h3>
    <el-button @click="load">刷新</el-button>
    <pre>{{ items }}</pre>
  </div>
</template>

<script setup lang="ts">
// 设置页当前仅提供列表查看，没有完整表单编辑器。
import { onMounted, ref } from 'vue'
import http from '../api/http'

const items = ref([] as any[])

const load = async () => {
  items.value = (await http.get('/api/settings')).data.data
}

onMounted(load)
</script>
