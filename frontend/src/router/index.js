import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
    meta: { title: '首页', requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login.vue'),
    meta: { title: '登录', requiresAuth: false, guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/register.vue'),
    meta: { title: '注册', requiresAuth: false, guestOnly: true }
  },
  {
    path: '/models',
    name: 'Models',
    component: () => import('@/views/models/index.vue'),
    meta: { title: '模型列表', requiresAuth: false }
  },
  {
    path: '/docs',
    name: 'Docs',
    component: () => import('@/views/docs/index.vue'),
    meta: { title: '接入说明', requiresAuth: false }
  },
  {
    path: '/user',
    component: () => import('@/layouts/UserLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/user/dashboard'
      },
      {
        path: 'dashboard',
        name: 'UserDashboard',
        component: () => import('@/views/user/dashboard.vue'),
        meta: { title: '控制台' }
      },
      {
        path: 'recharge',
        name: 'UserRecharge',
        component: () => import('@/views/user/recharge.vue'),
        meta: { title: '充值中心' }
      },
      {
        path: 'api-keys',
        name: 'UserAPIKeys',
        component: () => import('@/views/user/api-keys.vue'),
        meta: { title: 'API密钥管理' }
      },
      {
        path: 'call-logs',
        name: 'UserCallLogs',
        component: () => import('@/views/user/call-logs.vue'),
        meta: { title: '调用日志' }
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/profile.vue'),
        meta: { title: '个人信息' }
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/dashboard.vue'),
        meta: { title: '数据概览' }
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/users.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'providers',
        name: 'AdminProviders',
        component: () => import('@/views/admin/providers.vue'),
        meta: { title: '服务商管理' }
      },
      {
        path: 'models',
        name: 'AdminModels',
        component: () => import('@/views/admin/models.vue'),
        meta: { title: '模型管理' }
      },
      {
        path: 'payment-methods',
        name: 'AdminPaymentMethods',
        component: () => import('@/views/admin/payment-methods.vue'),
        meta: { title: '支付方式管理' }
      },
      {
        path: 'recharge-records',
        name: 'AdminRechargeRecords',
        component: () => import('@/views/admin/recharge-records.vue'),
        meta: { title: '充值记录' }
      },
      {
        path: 'call-logs',
        name: 'AdminCallLogs',
        component: () => import('@/views/admin/call-logs.vue'),
        meta: { title: '调用日志' }
      },
      {
        path: 'points-config',
        name: 'AdminPointsConfig',
        component: () => import('@/views/admin/points-config.vue'),
        meta: { title: '积分配置' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - Token Hub` : 'Token Hub - AI大模型中转平台'

  const userStore = useUserStore()
  const isAuthenticated = userStore.isAuthenticated

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.guestOnly && isAuthenticated) {
    next({ name: 'Home' })
    return
  }

  if (to.meta.requiresAdmin) {
    if (!userStore.userInfo && userStore.token) {
      await userStore.fetchUserInfo()
    }
    if (!userStore.isAdmin) {
      next({ name: 'Home' })
      return
    }
  }

  next()
})

export default router
