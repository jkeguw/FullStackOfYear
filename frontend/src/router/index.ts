// frontend/src/router/index.ts
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: () => import('@/views/Home.vue'),
        meta: { title: 'Home' }
    },
    {
        path: '/reviews',
        component: () => import('@/views/Reviews.vue'),
        meta: { title: 'Reviews' }
    },
    {
        path: '/tools',
        component: () => import('@/views/Tools.vue'),
        meta: { title: 'Tools' }
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;