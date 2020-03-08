<template>
  <div :class="$style.tag" v-on="$listeners">
    {{ label }}
  </div>
</template>

<script>
export default {
  props: {
    tag: {
      type: [Object, String],
    },
  },
  computed: {
    isObject () {
      return typeof this.tag == 'object'
    },
    label () {
      if (this.isObject) {
        return this.tag.name
      }
      return this.tag
    },
  },
}
</script>

<style lang="scss" module>
.tag {
  @apply text-xs;
  @apply inline-block;
  @apply px-2;
  @apply rounded;
  @apply bg-tag-bg;
  @apply whitespace-no-wrap;
}
</style>
