<template>
  <v-app-bar app clipped-left color="primary" dark>
    <v-app-bar-nav-icon @click="toggleDrawer"></v-app-bar-nav-icon>
    <span class="title pr-2" :style="titleStyles" v-if="!$vuetify.breakpoint.mobile">Pinmonl</span>

    <app-search></app-search>

    <v-spacer></v-spacer>
    <v-menu>
      <template #activator="{ on, attrs }">
        <v-btn icon v-on="on" v-bind="attrs">
          <v-icon>{{ mdiDotsVertical }}</v-icon>
        </v-btn>
      </template>
      <v-list dense>
        <v-list-item @click="toggleDarkMode">
          <v-list-item-content>
            <v-list-item-title>Dark mode</v-list-item-title>
          </v-list-item-content>
          <v-list-item-icon>
            <v-icon>{{ mdiMoonWaxingCrescent }}</v-icon>
          </v-list-item-icon>
        </v-list-item>

        <v-list-item @click="logout">
          <v-list-item-content>
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item-content>
          <v-list-item-icon>
            <v-icon>{{ mdiLogout }}</v-icon>
          </v-list-item-icon>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script>
import {
  mdiDotsVertical,
  mdiLogout,
  mdiMoonWaxingCrescent,
  mdiTagMultiple,
  mdiPackage,
} from '@mdi/js'
import * as auth from '@/utils/auth'
import AppSearch from './AppSearch'

export default {
  name: 'app-header',
  components: {
    AppSearch,
  },
  props: ['drawer'],
  data: () => ({
    mdiDotsVertical,
    mdiLogout,
    mdiMoonWaxingCrescent,
    mdiTagMultiple,
    mdiPackage,
  }),
  computed: {
    titleStyles () {
      const { mdAndUp } = this.$vuetify.breakpoint

      if (mdAndUp) {
        return { width: '204px' }
      }
      return {}
    },
  },
  methods: {
    toggleDrawer () {
      this.$emit('update:drawer', !this.drawer)
    },
    logout () {
      auth.logout()
    },
    toggleDarkMode () {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark
    },
  },
}
</script>
