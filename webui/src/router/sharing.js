import { abstractView } from './utils'
import Sharing from '@/views/Sharing.vue'

export default {
  path: '/sharing/:user/:name',
  component: abstractView,
  children: [
    {
      path: '',
      name: 'sharing',
      component: Sharing,
      props: true,
    },
    {
      path: 't/:parentName+',
      name: 'sharing.tag',
      component: Sharing,
      props: true,
    },
    {
      path: 'p/:pinlId',
      name: 'sharing.pinl',
      component: Sharing,
      props: true,
    },
  ],
}
