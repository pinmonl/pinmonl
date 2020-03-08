<template>
  <div :class="$style.pkg">
    <component :is="vendorComponent" :pkg="pkg" />
  </div>
</template>

<script>
import GithubVendor from './vendor/Github.vue'

export default {
  props: {
    pkg: {
      type: Object,
      required: true,
    },
  },
  computed: {
    vendor () {
      return this.pkg.vendor
    },
    vendorComponent () {
      const cMap = {
        'github': GithubVendor,
      }
      return cMap[this.vendor]
    },
  },
}
</script>

<style lang="scss" module>
.pkg {
}
</style>
