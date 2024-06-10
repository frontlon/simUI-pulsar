import {RouteRecordRaw} from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: () => import('pages/classic/Layout.vue'),
        meta: {
            refreshPage: true
        }
    },
    {
        path: '/playnite',
        component: () => import('pages/playnite/Layout.vue'),
        meta: {
            refreshPage: true
        }
    },
    {
        path: '/tiny',
        component: () => import('pages/tiny/Layout.vue'),
        meta: {
            refreshPage: true
        }
    },
    {
        path: '/config',
        component: () => import('pages/config/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/platform',
        component: () => import('pages/configPlatform/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/classic/ui',
        component: () => import('pages/classic/configUI/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/playnite/ui',
        component: () => import('pages/playnite/configUI/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/tiny/ui',
        component: () => import('pages/tiny/configUI/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/romManage',
        component: () => import('pages/romManage/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    {
        path: '/test',
        component: () => import('pages/test/Layout.vue'),
        meta: {
            refreshPage: true
        },
    },
    // Always leave this as last one,
    // but you can also remove it
    {
        path: '/:catchAll(.*)*',
        component: () => import('pages/ErrorNotFound.vue'),
    },
];

export default routes;
