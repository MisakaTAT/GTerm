import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHashHistory } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    redirect: 'content/connection',
    component: () => import('@/layouts/Layout.vue'),
    children: [
      {
        path: '/terminal',
        name: 'Terminal',
        component: () => import('@/views/pages/Terminal.vue'),
      },
      {
        path: 'content',
        name: 'Content',
        children: [
          {
            path: 'connection',
            name: 'Connection',
            component: () => import('@/views/pages/Connection.vue'),
          },
          {
            path: 'credential',
            name: 'Credential',
            component: () => import('@/views/pages/Credential.vue'),
          },
          {
            path: 'preferences',
            name: 'Preferences',
            component: () => import('@/views/pages/Preferences.vue'),
          },
        ],
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
