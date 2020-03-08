<template>
  <div :class="containerClass" v-on="$listeners">
    <svg :class="$style.icon">
      <use :xlink:href="`#${name}`" />
    </svg>
  </div>
</template>

<script>
export default {
  props: {
    name: {
      type: String,
    },
    size: {
      type: [Number, String],
      default: 1,
    },
  },
  computed: {
    containerClass () {
      return [{
        [this.$style.size_1x]: this.size === 1,
        [this.$style.size_fill]: this.size === 'fill',
      }]
    },
  },
}
</script>

<style lang="scss" module>
$size: 24px;

.size_1x {
  width: $size;
  height: $size;
}

.size_fill {
  @apply w-full;
  @apply h-full;
}

.icon {
  @apply block;
  @apply fill-current;
  @apply w-full;
  @apply h-full;
}
</style>
