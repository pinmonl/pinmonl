<template>
  <div :class="backdropClass" v-on="$listeners"></div>
</template>

<script>
export default {
  props: {
    absolute: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    backdropClass () {
      return [this.$style.backdrop, {
        [this.$style.backdrop_absolute]: this.absolute,
      }]
    },
  },
}
</script>

<style lang="scss" module>
.backdrop {
  @apply fixed;
  @apply inset-0;
}

.backdrop_absolute {
  @apply absolute;
}
</style>
