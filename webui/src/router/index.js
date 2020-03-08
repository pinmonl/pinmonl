import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import { abstractView } from './utils'
import Account from '@/views/Account.vue'
import Bookmark from '@/views/Bookmark.vue'
import Login from '@/views/Login.vue'
import Share from '@/views/Share.vue'
import Tag from '@/views/Tag.vue'
import routesToSharing from './sharing'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    redirect: { name: 'bookmark' },
    meta: { auth: true, showNav: true },
    component: abstractView,
    children: [
      {
        path: 'bookmark/new',
        name: 'bookmark.new',
        component: Bookmark,
        props: () => ({ new: true }),
      },
      {
        path: 'bookmark/:id?',
        name: 'bookmark',
        component: Bookmark,
        props: true,
      },
      {
        path: 'tag/new',
        name: 'tag.new',
        component: Tag,
        props: () => ({ new: true }),
      },
      {
        path: 'tag/:id?',
        name: 'tag',
        component: Tag,
        props: true,
      },
      {
        path: 'tag/c/:parentName+',
        name: 'tag.children',
        component: Tag,
        props: true,
      },
      {
        path: 'share/new',
        name: 'share.new',
        component: Share,
        props: { new: true },
      },
      {
        path: 'share/:id?',
        name: 'share',
        component: Share,
        props: true,
      },
      {
        path: 'account',
        name: 'account',
        component: Account,
      },
    ],
  },
  routesToSharing,
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { auth: false },
  },
  {
    path: '/signup',
    name: 'signup',
    component: Login,
    meta: { auth: false },
    props: () => ({ new: true }),
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

let appIsReady = false

router.beforeEach(async (to, from, next) => {
  if (!appIsReady) {
    const promise = new Promise(resolve => {
      const sleep = () => setTimeout(() => {
        if (!store.state['ready']) {
          sleep()
        } else {
          resolve()
        }
      }, 200)
      sleep()
    })
    await promise
    appIsReady = true
  }

  const authed = store.getters['authed']
  const metaAuth = to.matched.filter(m => typeof m.meta.auth == 'boolean')
  if (metaAuth.length > 0) {
    const needsAuth = metaAuth.some(m => m.meta.auth)
    if (!needsAuth && authed) {
      next({ name: 'home' })
      return
    } else if (needsAuth && !authed) {
      next({ name: 'login' })
      return
    }
  }

  next()
})

export default router
