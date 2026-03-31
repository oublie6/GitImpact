<template><div><h3>任务列表</h3><el-button @click="load">刷新</el-button><pre>{{ tasks }}</pre></div></template>
<script setup lang="ts">import {ref,onMounted} from 'vue';import http from '../api/http';const tasks=ref([] as any[]);const load=async()=>{tasks.value=(await http.get('/api/tasks')).data.data};onMounted(load)</script>
