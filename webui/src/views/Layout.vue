<template>
  <v-app>
    <app-header :drawer.sync="drawer"></app-header>
    <app-navbar
      :drawer.sync="drawer"
      @open-tag-editor="isTagEditorActive = true"
    ></app-navbar>

    <v-main>
      <v-container :class="containerClasses" fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import AppHeader from '@/components/app/AppHeader'
import AppNavbar from '@/components/app/AppNavbar'

export default {
  name: 'layout',
  components: {
    AppHeader,
    AppNavbar,
  },
  data: () => ({
    drawer: null,
  }),
  computed: {
    containerClasses () {
      return {
      }
    },
  },
}
</script>
