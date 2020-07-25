import Vue from 'vue'
import router from '@/router'
import store from '@/store'
import App from '@/App'
import vuetify from './plugins/vuetify';
// import Pubsub from '@/utils/pubsub'

new Vue({
  store,
  router,
  vuetify,
  render: h => h(App),
}).$mount('#app')

// new Pubsub(store.state.token)
