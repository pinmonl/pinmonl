<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import * as auth from '@/utils/auth'

export default {
  name: 'app',
  created () {
    this.refreshToken()
  },
  computed: {
    ...mapState(['client']),
  },
  methods: {
    async refreshToken () {
      try {
        const res = await this.client.refresh()
        auth.save(res)
      } catch (e) {
        //
      }
    },
  },
}
</script>

<style lang="scss">
@import './scss/styles';
</style>
