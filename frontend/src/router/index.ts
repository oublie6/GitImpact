import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import DashboardView from '../views/DashboardView.vue'
import RepositoriesView from '../views/RepositoriesView.vue'
import TaskCreateView from '../views/TaskCreateView.vue'
import TaskListView from '../views/TaskListView.vue'
import TaskDetailView from '../views/TaskDetailView.vue'
import SettingsView from '../views/SettingsView.vue'

const router = createRouter({ history: createWebHistory(), routes: [
  { path: '/login', component: LoginView },
  { path: '/register', component: RegisterView },
  { path: '/', component: DashboardView },
  { path: '/repositories', component: RepositoriesView },
  { path: '/tasks/new', component: TaskCreateView },
  { path: '/tasks', component: TaskListView },
  { path: '/tasks/:id', component: TaskDetailView },
  { path: '/settings', component: SettingsView }
]})
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!auth.token && !['/login', '/register'].includes(to.path)) return '/login'
  if (auth.token && to.path === '/login') return '/'
})
export default router
