<template>
  <div :class="$style.group">
    <slot></slot>
    <div :class="$style.errors">
      <slot name="errors"></slot>
    </div>
  </div>
</template>

<style lang="scss" module>
.group {
}

.errors {
  @apply mb-4;
  @apply text-xs;
  @apply text-error;
  @apply px-1;
}
</style>