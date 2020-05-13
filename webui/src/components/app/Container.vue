<template>
  <div :class="$style.container">
    <Header :class="[$style.header, headerClass]">
      <slot name="header"></slot>
    </Header>
    <div :class="bodyClasses">
      <slot></slot>
    </div>
  </div>
</template>

<script>
import Header from './Header.vue'

export default {
  components: {
    Header,
  },
  props: {
    headerClass: {
      type: [String, Array],
      default: null,
    },
    bodyClass: {
      type: [String, Array],
      default: null,
    },
    overflowY: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    bodyClasses () {
      return [this.$style.body, this.bodyClass, {
        [this.$style.body_overflowY]: this.overflowY,
      }]
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply relative;
  padding-top: theme('height.header');
  @apply h-full;
}

.header {
  @apply bg-container;
  @apply fixed;
  @apply z-100;
  @apply inset-x-0;
  @apply top-0;
}

.body {
  @apply relative;
  @apply h-full;
  @apply w-full;
}

.body_overflowY {
  @apply overflow-y-auto;
  @apply scrolling-touch;
}
</style>
