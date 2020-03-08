<template>
  <div :class="$style.container">
    <slot></slot>
  </div>
</template>

<style lang="scss" module>
.container {
  @apply mb-2;
  @apply mx-3;
  @apply text-xs;
  @apply text-error;

  > p {
    @apply mb-1;
  }
}
</style>