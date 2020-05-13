<template>
  <div :class="$baseStyle.container">
    <Stat :class="$baseStyle.stat" :value="pkg.providerUri" />
    <Stat :class="$baseStyle.stat" kind="tag" :value="latestTag" title="Latest Tag" />
    <Stat :class="$baseStyle.stat" kind="star" :value="starCount" as-number title="Star" />
    <Stat :class="$baseStyle.stat" kind="issue" :value="issueCount" as-number title="Open Issue" />
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
    issueCount () {
      return this.findStatValue('open_issue')
    },
  },
}
</script>
