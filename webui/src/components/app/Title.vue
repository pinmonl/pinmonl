<template>
  <h2 :class="$style.title">
    <slot></slot>
  </h2>
</template>

<style lang="scss" module>
.title {
  @apply mb-2;
}
</style>
