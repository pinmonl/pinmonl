import Vue from 'vue'
import router from '@/router'
import store from '@/store'
import App from '@/App'
import vuetify from './plugins/vuetify';

new Vue({
  store,
  router,
  vuetify,
  render: h => h(App),
}).$mount('#app')
