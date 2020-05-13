<template>
  <div :class="$baseStyle.container">
    <Stat :class="$baseStyle.stat" :value="pkg.providerUri" />
    <Stat :class="$baseStyle.stat" kind="tag" :value="latestTag" title="Latest Tag" />
    <Stat :class="$baseStyle.stat" kind="size" :value="fileSize" as-file-size title="Package Size" />
    <Stat :class="$baseStyle.stat" kind="download" :value="downloadCount" as-large-number title="Download" />
  </div>
</template>

<script>
import mixin from './mixin'

export default {
  mixins: [mixin],
  computed: {
    fileSize () {
      return this.findStatValue('size')
    },
    downloadCount() {
      return this.findStatValue('download')
    },
  },
}
</script>
