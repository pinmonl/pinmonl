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
  @apply text-gray-600;
  @apply text-sm;
  @apply block;
  @apply p-1;
}
</style>
