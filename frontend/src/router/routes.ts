import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { title: '登录', requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/AdminLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' },
      },
      {
        path: 'tenants',
        name: 'Tenants',
        component: () => import('@/views/tenants/TenantListView.vue'),
        meta: { title: '租户管理', icon: 'OfficeBuilding' },
      },
      {
        path: 'tenants/:id',
        name: 'TenantDetail',
        component: () => import('@/views/tenants/TenantDetailView.vue'),
        meta: { title: '租户详情', hidden: true },
      },
      {
        path: 'model-categories',
        name: 'ModelCategories',
        component: () => import('@/views/model-categories/ModelCategoryListView.vue'),
        meta: { title: '物模型分类', icon: 'Collection' },
      },
      {
        path: 'thing-models',
        name: 'ThingModels',
        component: () => import('@/views/thing-models/ThingModelListView.vue'),
        meta: { title: '物模型管理', icon: 'Setting' },
      },
      {
        path: 'thing-models/:id',
        name: 'ThingModelDetail',
        component: () => import('@/views/thing-models/ThingModelDetailView.vue'),
        meta: { title: '物模型详情', hidden: true },
      },
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('@/views/devices/DeviceListView.vue'),
        meta: { title: '设备管理', icon: 'Monitor' },
      },
      {
        path: 'devices/:id',
        name: 'DeviceDetail',
        component: () => import('@/views/devices/DeviceDetailView.vue'),
        meta: { title: '设备详情', hidden: true },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404/NotFoundView.vue'),
    meta: { title: '404', requiresAuth: false },
  },
]
