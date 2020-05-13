<template>
  <Icon :name="kindIcon" :class="$style.icon" size="auto" />
</template>

<script>
export default {
  props: {
    kind: {
      type: String,
    },
  },
  computed: {
    kindIcon () {
      const kind = this.kind
      return {
        'issue': 'issue',
        'size':  'folderOutline',
        'star':  'star',
        'tag':   'tag',
        'download': 'download',
      }[kind]
    },
  },
}
</script>

<style lang="scss" module>
$size: 16px;

.icon {
  width: $size;
  height: $size;
}
</style>
