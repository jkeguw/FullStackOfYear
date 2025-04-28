import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 测试页面
    {
      path: '/test-comparison',
      name: 'TestComparison',
      component: () => import('../pages/TestComparisonPage.vue')
    },
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
      path: '/similar',
      name: 'SimilarMouse',
      component: () => import('../pages/SimilarMousePage.vue')
    },
    {
      path: '/tools',
      name: 'Tools',
      component: () => import('../pages/ToolsPage.vue')
    },
    {
      path: '/tools/ruler',
      name: 'RulerTool',
      component: () => import('../pages/RulerToolPage.vue')
    },
    {
      path: '/tools/sensitivity',
      name: 'SensitivityTool',
      component: () => import('../pages/SensitivityToolPage.vue')
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
    // 购物车路由
    {
      path: '/cart',
      name: 'Cart',
      component: () => import('../pages/CartPage.vue'),
      meta: { requiresAuth: true }
    },
    // 关于页面路由
    {
      path: '/about',
      name: 'About',
      component: () => import('../pages/AboutPage.vue')
    },
    {
      path: '/contact',
      name: 'Contact',
      component: () => import('../pages/ContactPage.vue')
    },
    {
      path: '/privacy',
      name: 'Privacy',
      component: () => import('../pages/PrivacyPage.vue')
    },
    {
      path: '/terms',
      name: 'Terms',
      component: () => import('../pages/TermsPage.vue')
    },
    // 捕获所有未定义路由，重定向到首页
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  }
});

// 全局导航守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);

  // 如果需要登录但用户未登录，重定向到登录页
  if (requiresAuth && !userStore.token) {
    next({ name: 'Login', query: { redirect: to.fullPath } });
  } else {
    next();
  }
});

export default router;
