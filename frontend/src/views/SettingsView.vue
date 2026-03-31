<template><div><h3>系统设置</h3><el-button @click="load">刷新</el-button><pre>{{ items }}</pre></div></template>
<script setup lang="ts">import {ref,onMounted} from 'vue';import http from '../api/http';const items=ref([] as any[]);const load=async()=>{items.value=(await http.get('/api/settings')).data.data};onMounted(load)</script>
