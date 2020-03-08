import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import globalMixin from './mixins/global'

Vue.mixin(globalMixin)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
