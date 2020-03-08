<template>
  <Vendor name="github" :url="githubURL" v-if="isReady">
    <Stat icon="cubeOutline" :value="vendorUri" />
    <Stat icon="star" :value="starStat.value" number />
    <Stat icon="tag" :value="tagStat.value" />
    <Stat icon="sourceFork" :value="forkStat.value" number />
    <Stat icon="eye" :value="watchStat.value" number />
  </Vendor>
</template>

<script>
import vendorMixin from './mixin'

export default {
  mixins: [vendorMixin],
  computed: {
    forkStat () {
      return this.stats.fork
    },
    watchStat () {
      return this.stats.watcher
    },
    starStat () {
      return this.stats.star
    },
    tagStat () {
      return this.stats.tag
    },
    githubURL () {
      return `https://github.com/${this.vendorUri}`
    },
  },
}
</script>
