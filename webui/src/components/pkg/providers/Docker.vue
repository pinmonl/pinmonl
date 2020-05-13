<template>
  <div :class="$baseStyle.container">
    <Stat :class="$baseStyle.stat" :value="pkg.providerUri" />
    <Stat :class="$baseStyle.stat" kind="tag" :value="latestTag" title="Latest Tag" />
    <Stat :class="$baseStyle.stat" kind="star" :value="starCount" title="Star" as-number />
    <Stat :class="$baseStyle.stat" kind="download" :value="pullCount" as-large-number title="Pull" />
  </div>
</template>

<script>
import mixin from './mixin'

export default {
  mixins: [mixin],
  computed: {
    starCount () {
      return this.findStatValue('star')
    },
    pullCount () {
      return this.findStatValue('pull')
    },
  },
}
</script>
