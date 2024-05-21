import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
  },
  {
    path: '/tests',
    name: 'Tests',
    component: () => import('../views/TestsView.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/RegistrationView.vue')
  },
  {
    path: "/accounts",
    name: "Accounts",
    component: () => import("../views/AccountsView.vue"),
    meta: {
      requireLogin: true,
    },
  },
  {
    path: "/tests_admin",
    name: "TestsAdmin",
    component: () => import("../views/AdminTestsView.vue"),
    meta: {
      requireLogin: true,
    },
  },
  {
    path: "/results",
    name: "Results",
    component: () => import("../views/AccountsView.vue"),
    meta: {
      requireLogin: true,
    },
  },
  {
    path: "/test/:testId",
    name: "TestFull",
    component: () => import("../views/TestView.vue"),
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
