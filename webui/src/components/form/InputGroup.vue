<template>
  <div :class="$style.group">
    <slot></slot>
    <slot name="errors">
      <div :class="$style.errors"></div>
    </slot>
  </div>
</template>

<style lang="scss" module>
.group {
}

.errors {
  @apply mb-4;
}
</style>