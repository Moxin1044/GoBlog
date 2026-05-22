import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import FrontLayout from '@/layouts/FrontLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: FrontLayout,
    children: [
      { path: '', name: 'Home', component: () => import('@/views/front/Home.vue') },
      { path: 'article/:id', name: 'ArticleDetail', component: () => import('@/views/front/ArticleDetail.vue') },
      { path: 'category/:id', name: 'CategoryPage', component: () => import('@/views/front/CategoryPage.vue') },
      { path: 'login', name: 'Login', component: () => import('@/views/front/Login.vue') },
      { path: 'register', name: 'Register', component: () => import('@/views/front/Register.vue') },
      { path: 'profile', name: 'Profile', component: () => import('@/views/front/Profile.vue'), meta: { requiresAuth: true } },
    ],
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAdmin: true },
    children: [
      { path: '', name: 'Dashboard', component: () => import('@/views/admin/Dashboard.vue') },
      { path: 'articles', name: 'ArticleManage', component: () => import('@/views/admin/ArticleManage.vue') },
      { path: 'article/create', name: 'ArticleCreate', component: () => import('@/views/admin/ArticleEdit.vue') },
      { path: 'article/edit/:id', name: 'ArticleEdit', component: () => import('@/views/admin/ArticleEdit.vue') },
      { path: 'categories', name: 'CategoryManage', component: () => import('@/views/admin/CategoryManage.vue') },
      { path: 'tags', name: 'TagManage', component: () => import('@/views/admin/TagManage.vue') },
      { path: 'navigations', name: 'NavigationManage', component: () => import('@/views/admin/NavigationManage.vue') },
      { path: 'comments', name: 'CommentManage', component: () => import('@/views/admin/CommentManage.vue') },
      { path: 'users', name: 'UserManage', component: () => import('@/views/admin/UserManage.vue') },
      { path: 'admins', name: 'AdminManage', component: () => import('@/views/admin/AdminManage.vue'), meta: { superAdmin: true } },
      { path: 'config', name: 'SystemConfig', component: () => import('@/views/admin/SystemConfig.vue'), meta: { superAdmin: true } },
      { path: 'ai-models', name: 'AIModelManage', component: () => import('@/views/admin/AIModelManage.vue'), meta: { superAdmin: true } },
      { path: 'backups', name: 'BackupManage', component: () => import('@/views/admin/BackupManage.vue'), meta: { superAdmin: true } },
      { path: 'logs', name: 'OperationLog', component: () => import('@/views/admin/OperationLog.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('role')

  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresAdmin) {
    if (!token || (userRole !== 'admin' && userRole !== 'superadmin')) {
      next({ name: 'Login' })
    } else if (to.meta.superAdmin && userRole !== 'superadmin') {
      next({ name: 'Dashboard' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
