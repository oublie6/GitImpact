<template><div><h3>仓库管理</h3><el-button @click="load">刷新</el-button><pre>{{ repos }}</pre></div></template>
<script setup lang="ts">import {ref,onMounted} from 'vue';import http from '../api/http';const repos=ref([] as any[]);const load=async()=>{repos.value=(await http.get('/api/repositories')).data.data};onMounted(load)</script>
