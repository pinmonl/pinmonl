<template>
  <button
    :class="buttonClass"
    v-bind="buttonProps"
    v-on="$listeners"
  >
    <slot></slot>
  </button>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'button',
    },
    noStyle: {
      type: Boolean,
      default: false,
    },
    block: {
      type: Boolean,
      default: false,
    },
    light: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    buttonProps () {
      return {
        type: this.type,
      }
    },
    buttonClass () {
      return [{
        [this.$style.button]: !this.noStyle,
        [this.$style.buttonBlock]: this.block,
        [this.$style.buttonLight]: this.light,
      }]
    }
  },
}
</script>

<style lang="scss" module>
.button {
  @apply bg-primary;
  @apply text-text-inverted;
  @apply px-4;
  @apply py-1;
  @apply rounded;
}

.buttonBlock {
  @apply block;
}

.buttonLight {
  @apply bg-btn-light-bg;
  @apply text-btn-light;
}
</style>
