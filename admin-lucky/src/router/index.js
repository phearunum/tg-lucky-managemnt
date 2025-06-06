import { createWebHistory, createRouter } from 'vue-router'
import Layout from '@/layout'

export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue')
      }]
  },
  {
    path: '/login',
    component: () => import('@/views/login'),
    hidden: true
  },
  {
    path: '/register',
    component: () => import('@/views/register'),
    hidden: true
  },
  {
    path: "/:pathMatch(.*)*",
    component: () => import('@/views/error/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error/401'),
    hidden: true
  },
  {
    path: '',
    component: Layout,
    redirect: '/index',
    children: [
      {
        path: '/index',
        component: () => import('@/views/dashboard/index.vue'),
        name: 'Index',
        meta: { title: 'Home', icon: 'dashboard', affix: true, titleKey: 'menu.home' }
      }]
  },
  {
    path: '/user',
    component: Layout,
    hidden: true,
    redirect: 'noredirect',
    children: [
      {
        path: 'profile',
        component: () => import('@/views/system/user/profile/index'),
        name: 'Profile',
        meta: { title: 'Profile', icon: 'user', titleKey: 'menu.personalCenter' }
      }]
  },

  {
    path: '',
    component: Layout,
    hidden: true,
    meta: { title: 'Component example', icon: 'icon', noCache: 'fasle' },

    children: [
      {
        path: 'icon',
        component: () => import('@/views/components/icons/index'),
        name: 'icon',
        meta: { title: 'icon', icon: 'icon1', noCache: 'fasle', titleKey: 'menu.icon' }
      }]
  },

];

const router = createRouter({
  mode: 'history',
  history: createWebHistory(
    import.meta.env.VITE_APP_ROUTER_PREFIX),
  routes: constantRoutes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
});

export default router;
