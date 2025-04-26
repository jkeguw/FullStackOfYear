import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('../pages/Home.vue')
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../pages/LoginPage.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../pages/RegisterPage.vue')
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../pages/ProfilePage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/devices',
      name: 'DeviceList',
      component: () => import('../pages/DeviceListPage.vue')
    },
    {
      path: '/devices/:id',
      name: 'DeviceDetail',
      component: () => import('../pages/DeviceDetailPage.vue')
    },
    {
      path: '/mice/:id',
      name: 'MouseDetail',
      component: () => import('../pages/MouseDetailPage.vue')
    },
    {
      path: '/compare',
      name: 'Compare',
      component: () => import('../pages/ComparePage.vue')
    },
    {
      path: '/tools',
      name: 'Tools',
      component: () => import('../pages/ToolsPage.vue'),
      children: [
        {
          path: '',
          name: 'ToolsList',
          component: () => import('../components/tools/ToolsList.vue')
        },
        {
          path: 'dpi',
          name: 'DpiCalculator',
          component: () => import('../components/tools/DPICalculator.vue')
        },
        {
          path: 'ruler',
          name: 'RulerTool',
          component: () => import('../pages/RulerToolPage.vue')
        },
        {
          path: 'sensitivity',
          name: 'SensitivityTool',
          component: () => import('../pages/SensitivityToolPage.vue')
        }
      ]
    },
    // 已移除个人设备管理功能
    
    {
      path: '/i18n-demo',
      name: 'I18nDemo',
      component: () => import('../pages/I18nDemoPage.vue')
    },
    {
      path: '/reviews',
      name: 'ReviewList',
      component: () => import('../pages/ReviewListPage.vue')
    },
    {
      path: '/reviews/:id',
      name: 'ReviewDetail',
      component: () => import('../pages/ReviewDetailPage.vue')
    },
    {
      path: '/reviews/create',
      name: 'ReviewCreate',
      component: () => import('../pages/ReviewForm.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/database',
      name: 'MouseDatabase',
      component: () => import('../pages/MouseDatabasePage.vue')
    },
    // 订单相关路由
    {
      path: '/orders',
      name: 'OrderList',
      component: () => import('../pages/OrderListPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/orders/:id',
      name: 'OrderDetail',
      component: () => import('../pages/OrderDetailPage.vue'),
      meta: { requiresAuth: true }
    },
    // 结账流程路由
    {
      path: '/checkout',
      name: 'Checkout',
      component: () => import('../pages/CheckoutPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/checkout/payment/:id',
      name: 'OrderPayment',
      component: () => import('../pages/CheckoutPage.vue'),
      props: (route) => ({ 
        orderId: route.params.id,
        step: 'payment'
      }),
      meta: { requiresAuth: true }
    },
    // 购物车路由
    {
      path: '/cart',
      name: 'Cart',
      component: () => import('../pages/CartPage.vue')
    },
  ]
})

// 全局导航守卫
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isLoggedIn = !!localStorage.getItem('token')

  if (requiresAuth && !isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router