<template>
  <v-chip small outlined :class="pkgClasses" v-if="provider">
    <v-avatar :left="!iconOnly">
      <v-icon small>{{ provider.icon }}</v-icon>
    </v-avatar>
    <template v-if="!iconOnly">
      <span>{{ pkg.providerUri }}</span>
      <template v-if="latestTag">
        <span>:</span>
        <span>{{ latestTag }}</span>
      </template>
    </template>
  </v-chip>
</template>

<script>
import { getProvider } from '@/utils/pkg'

export default {
  name: 'pkg-chip',
  props: ['pkg'],
  computed: {
    provider () {
      return getProvider(this.pkg.provider)
    },
    latestTag () {
      const tag = this.pkg.stats.find(stat => stat.kind == 'tag')
      if (!tag || !tag.isLatest) {
        return
      }
      return tag.value
    },
    iconOnly () {
      return !this.$vuetify.breakpoint.mdAndUp
    },
    pkgClasses () {
      return {
        [this.$style.pkgSm]: this.$vuetify.breakpoint.smAndDown,
      }
    },
  },
  methods: {
  },
}
</script>

<style lang="scss" module>
.pkgSm {
  padding: 0;
}
</style>
