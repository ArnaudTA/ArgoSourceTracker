import {
    createRouter,
    createWebHistory,
} from 'vue-router'
import type {
    RouteRecordRaw,
} from 'vue-router'


const Home = () => import('../views/Home.vue')
const NotFound = () => import('../views/NotFound.vue')
const ApplicationDetails = () => import('../views/ApplicationDetails.vue')

const MAIN_TITLE = 'ArgoCD Source Tracker'

export const routes: Readonly<RouteRecordRaw[]> = [
    {
        path: '/',
        redirect: { name: 'Home' }
    }, {
        path: '/ui',
        children: [
            {
                path: '',
                name: 'Home',
                component: Home,
            },
            {
                path: 'application/:name',
                name: 'Application',
                props(to) {
                    return {
                        name: to.params.name
                    }
                },
                component: ApplicationDetails,
            },
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: NotFound,
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env?.BASE_URL || ''),
    routes,
})

/**
 * Set application title
 */
router.beforeEach((to) => { // Cf. https://github.com/vueuse/head pour des transformations avancées de Head
    const specificTitle = to.meta.title ? `${to.meta.title} - ` : ''
    document.title = `${specificTitle}${MAIN_TITLE}`
})

export default router