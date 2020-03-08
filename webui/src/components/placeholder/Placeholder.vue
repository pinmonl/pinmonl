<template>
  <div :class="$style.placeholder" v-on="listeners">
    <slot></slot>
  </div>
</template>

<script>
export default {
  computed: {
    listeners () {
      return { ...this.$listeners }
    },
  },
}
</script>

<style lang="scss" module>
.placeholder {
  @apply relative;
  @apply overflow-hidden;

  &:before {
    content: '';
    @apply absolute;
    @apply inset-0;
    width: 1000px;
    mix-blend-mode: overlay;
    background: linear-gradient(90deg, #cccccc 34%, #666666 50%, #cccccc 66%);
    background-size: 80%;
    animation: shimmer 3.5s infinite linear;
  }

  > div {
    background-color: #d9d9d9;
  }
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}
</style>
