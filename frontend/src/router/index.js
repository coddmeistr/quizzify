import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/tests',
    name: 'Tests',
    component: () => import('../views/TestsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
