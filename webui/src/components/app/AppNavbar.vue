<template>
  <v-navigation-drawer app v-model="vDrawer" clipped>
    <v-list dense nav>
      <v-list-item
        v-for="item in items"
        :key="item.title"
        :title="item.title"
        :to="item.to"
      >
        <v-list-item-icon>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { mdiBookmark, mdiTag } from '@mdi/js'

export default {
  name: 'app-navbar',
  props: ['drawer'],
  computed: {
    vDrawer: {
      get () { return this.drawer },
      set (v) { this.$emit('update:drawer', v) }
    },
    items () {
      return [
        { title: 'Pins', icon: mdiBookmark, to: '/pin' },
        { title: 'Tags', icon: mdiTag, to: '/tag' },
        // { title: 'Shares', icon: mdiShareVariant, to: '/share' },
      ]
    }
  },
}
</script>
