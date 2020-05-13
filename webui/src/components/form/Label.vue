<template>
  <label
    :class="$style.label"
    v-on="$listeners"
  >
    <slot></slot>
  </label>
</template>

<style lang="scss" module>
.label {
  @apply text-text-secondary;
  @apply text-sm;
  @apply block;
  @apply p-1;
}
</style>
