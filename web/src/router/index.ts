import { toRefs } from 'vue'
import { until } from '@vueuse/core'
import { storeToRefs } from 'pinia'

import { createRouter, createWebHistory, type RouteMeta } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { usePocketbaseStore } from '@/stores/pb'
import { type RecordModel } from 'pocketbase'

import 'vue-router'
export {}

declare module 'vue-router' {
  interface RouteMeta {
    // is optional
    userRoles?: Array<string>
    // must be declared by every route
    requiresAuth?: boolean
  }
}
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      name: 'home',
      path: '/',
      component: HomeView
    },
    {
      name: 'admin-page',
      path: '/admin',
      component: () => import('../views/AdminHome.vue')
    },
    {
      name: 'session-create',
      path: '/admin/session',
      meta: { requiresAuth: true, userRoles: ['admin'] },
      component: () => import('../views/SessionCreate.vue')
    },
    {
      name: 'session-create',
      path: '/admin/session',
      meta: { requiresAuth: true, userRoles: ['admin'] },
      component: () => import('../views/SessionCreate.vue')
    },
    {
      name: 'session-overview',
      path: '/session/:sessionName',
      props: true,
      meta: { requiresAuth: true, userRoles: ['super', 'admin'] },
      component: () => import('../views/SessionView.vue'),
      beforeEnter: async (to, from) => {
        const pocketbaseStore = usePocketbaseStore()
        const { api, isLoggedIn, user } = pocketbaseStore
        if (isLoggedIn && user?.role === 'super') {
          const isPwdSet = await api.send('/auth-set', { method: 'GET' })
          if (isPwdSet === 'true' || isPwdSet === true) {
            return true
          } else {
            return {
              name: 'change-pwd',
              replace: true
            }
          }
        } else {
          return true
        }
      }
    },
    {
      name: 'change-pwd',
      path: '/change-password',
      props: true,
      meta: { requiresAuth: true, userRoles: ['super'] },
      component: () => import('../views/ChangePassword.vue')
    },
    {
      name: 'session-analysis2',
      path: '/session/:sessionName/analysis2',
      props: true,
      meta: { requiresAuth: true, userRoles: ['super', 'admin'] },
      component: () => import('../views/Analysis2View.vue')
    },
    {
      name: 'session-analysis3',
      path: '/session/:sessionName/analysis3',
      props: true,
      meta: { requiresAuth: true, userRoles: ['super', 'admin'] },
      component: () => import('../views/Analysis3View.vue')
    },
    {
      name: 'session-form-welcome',
      path: '/session/:sessionName/survey',
      props: true,
      meta: { requiresAuth: true },
      component: () => import('../views/SurveyWelcome.vue')
    },
    {
      name: 'session-form',
      path: '/session/:sessionName/survey/:username',
      meta: { requiresAuth: true, match: true },
      component: () => import('../views/SurveyWrapper.vue')
    },
    {
      name: 'session-form-end',
      path: '/session/:sessionName/survey/:username/end',
      props: true,
      meta: { requiresAuth: true },
      component: () => import('../views/SurveyEnd.vue')
    },
    {
      name: 'session-form-view',
      path: '/session/:sessionName/survey/:username/view',
      props: true,
      meta: { requiresAuth: true },
      component: () => import('../views/SurveySummaryView.vue')
    },
    {
      name: 'new-form',
      path: '/new-form',
      component: () => import('../views/NewForm.vue')
    }
  ]
})

function checkPrivilege(meta: RouteMeta, user?: RecordModel) {
  if (meta.userRoles) {
    if (meta.userRoles.includes('admin') && user?.collectionName == '_superusers') {
      return true
    }
    if (meta.userRoles.includes('super') && user?.role === 'super') {
      return true
    } else {
      return false
    }
  }
  return true
}

router.beforeEach(async (to, from, next) => {
  const pocketbaseStore = usePocketbaseStore()
  const { isLoggedIn, login, logout } = pocketbaseStore

  if (to.query['u'] && to.query['t']) {
    const username: string = to.query['u'] as string
    const token: string = to.query['t'] as string

    if (isLoggedIn) {
      await logout()
    }
    await login(username, token)

    delete to.query['u']
    delete to.query['t']
  }
  next()
})

router.beforeEach(async (to, from) => {
  const pocketbaseStore = usePocketbaseStore()
  const { isLoggedIn, user } = pocketbaseStore
  if (to.name === 'home') {
    if (isLoggedIn && user?.avatar !== undefined) {
      return {
        name: 'session-create',
        replace: true
      }
    } else if (isLoggedIn) {
      const { session } = storeToRefs(pocketbaseStore)
      await until(session).changed()
      if (isLoggedIn && user?.role && user.role === 'super') {
        return {
          name: 'session-overview',
          params: { sessionName: session.value?.name },
          replace: true
        }
      } else if (isLoggedIn && user?.role && user.role === 'normal') {
        return {
          name: 'session-form-welcome',
          params: { sessionName: session.value?.name },
          replace: true
        }
      }
    }

    return
  }
  if (to.meta.requiresAuth) {
    if (!isLoggedIn) {
      // this route requires auth, check if logged in
      // if not, redirect to login page.
      console.error('access restrict, user is not logged in')
      return {
        name: 'home',
        replace: true
        //) query: { redirect: to.fullPath }
      }
    } else if (!checkPrivilege(to.meta, user)) {
      console.error('access restrict, user do not have privilege.')
      return {
        name: 'home',
        replace: true
      }
    }

    if (to.meta.match) {
      if (to.params['username'] !== user?.username) {
        console.error('access restrict, user does not match requested content.')
        return {
          name: 'home',
          replace: true
        }
      }
    }
  }
})

export default router
