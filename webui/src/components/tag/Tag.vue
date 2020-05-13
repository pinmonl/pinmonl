<template>
  <div :class="tagClasses" v-on="$listeners">
    {{ label }}
  </div>
</template>

<script>
export default {
  props: {
    tag: {
      type: [Object, String],
    },
    lg: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    tagClasses () {
      return [this.$style.tag, {
        [this.$style.tag_lg]: this.lg,
      }]
    },
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
  @apply border;
  @apply border-tag-border;
}

.tag_lg {
  @apply text-sm;
}
</style>
