import Vue from 'vue'
import VueRouter from 'vue-router'
import Layout from '@/views/Layout'
import Pin from '@/views/Pin'
import Tag from '@/views/Tag'
import Login from '@/views/Login'
import { baseURL } from '@/utils/constants'
import store from '@/store'

Vue.use(VueRouter)

const router = new VueRouter({
  base: baseURL,
  mode: 'history',
  routes: [
    {
      path: '/login',
      component: Login,
      beforeEnter: (to, from, next) => {
        if (store.getters.isAuthed) {
          next('/')
          return
        }

        next()
      },
    },
    {
      path: '/',
      component: Layout,
      meta: {
        requireAuth: true,
      },
      children: [
        {
          path: 'pin',
          component: Pin,
        },
        {
          path: 'tag',
          component: Tag,
        },
        {
          path: '*',
          redirect: {
            path: '/pin',
          },
        },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(route => route.meta.requireAuth)) {
    if (!store.getters.isAuthed) {
      next('/login')
      return
    }
  }

  next()
})

export default router
