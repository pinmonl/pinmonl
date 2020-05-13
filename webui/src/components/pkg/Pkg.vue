<template>
  <div :class="$style.container">
    <component :is="providerComponent" :pkg="pkg" />
  </div>
</template>

<script>
import DockerProvider from './providers/Docker.vue'
import GithubProvider from './providers/Github.vue'
import HelmProvider from './providers/Helm.vue'
import NpmProvider from './providers/Npm.vue'

export default {
  props: {
    pkg: {
      type: Object,
    },
  },
  computed: {
    providerComponent () {
      return {
        'docker': DockerProvider,
        'github': GithubProvider,
        'helm':   HelmProvider,
        'npm':    NpmProvider,
      }[this.pkg.provider]
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply text-xs;
}
</style>
