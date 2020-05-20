import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import Account from '@/views/Account.vue'
import Bookmark from '@/views/Bookmark.vue'
import Container from '@/views/Container.vue'
import Login from '@/views/Login.vue'
import Share from '@/views/Share.vue'
import Sponsors from '@/views/Sponsors.vue'
import SupportUs from '@/views/SupportUs.vue'
import Tag from '@/views/Tag.vue'
import ThreeColumn from '@/views/ThreeColumn.vue'
import routesToSharing from './sharing'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    redirect: { name: 'bookmark.list' },
    meta: { auth: true, showNav: true },
    component: Container,
    children: [
      {
        path: 'bookmark',
        component: ThreeColumn,
        meta: { navShowNewBookmark: true },
        children: [
          {
            path: '',
            name: 'bookmark.list',
            component: Bookmark,
          },
          {
            path: 'new',
            name: 'bookmark.new',
            component: Bookmark,
            props: () => ({ isNew: true }),
            meta: { noSearch: true },
          },
          {
            path: ':id',
            name: 'bookmark',
            component: Bookmark,
            props: true,
            meta: { noSearch: true },
          },
        ],
      },
      {
        path: 'tag',
        component: ThreeColumn,
        meta: { navShowNewTag: true },
        children: [
          {
            path: '',
            name: 'tag.list',
            component: Tag,
          },
          {
            path: 'new',
            name: 'tag.new',
            component: Tag,
            props: () => ({ isNew: true }),
            meta: { noSearch: true },
          },
          {
            path: ':id',
            name: 'tag',
            component: Tag,
            props: true,
            meta: { noSearch: true },
          },
          {
            path: 'c/:parentName+',
            name: 'tag.children',
            component: Tag,
            props: true,
          },
        ],
      },
      {
        path: 'share/new',
        name: 'share.new',
        component: Share,
        props: () => ({ isNew: true }),
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
      {
        path: 'sponsors',
        name: 'sponsors',
        component: Sponsors,
      },
      {
        path: 'support-us',
        name: 'support-us',
        component: SupportUs,
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
    props: () => ({ isNew: true }),
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
