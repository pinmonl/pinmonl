<template>
  <div :class="$baseStyle.container">
    <Stat :class="$baseStyle.stat" :value="pkg.providerUri" />
    <Stat :class="$baseStyle.stat" kind="tag" :value="latestTag" title="Latest tag" />
  </div>
</template>

<script>
import mixin from './mixin'

export default {
  mixins: [mixin],
  computed: {
  },
}
</script>
