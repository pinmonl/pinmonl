<template>
  <button v-on="$listeners" :class="buttonClass">
    <Icon :name="name" :class="$style.icon" />
    <slot></slot>
  </button>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
      required: true,
    },
    block: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    buttonClass () {
      return [this.$style.button, {
        [this.$style.buttonBlock]: this.block,
      }]
    },
  },
}
</script>

<style lang="scss" module>
.button {
  width: 40px;
  height: 40px;

  .icon {
    @apply p-2;
    @apply h-full;
    @apply w-full;
  }
}

.buttonBlock {
  @apply block;
}
</style>
