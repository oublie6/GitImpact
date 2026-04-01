<template>
  <el-card>
    <el-form>
      <el-input v-model="username" placeholder="用户名" />
      <el-input v-model="email" placeholder="邮箱" />
      <el-input v-model="password" type="password" placeholder="密码" />
      <el-button @click="submit">注册</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
// 注册页：成功后回到登录页，由用户手动登录。
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const s = useAuthStore()
const r = useRouter()
const username = ref('')
const email = ref('')
const password = ref('')

const submit = async () => {
  await s.register({ username: username.value, email: email.value, password: password.value })
  r.push('/login')
}
</script>
