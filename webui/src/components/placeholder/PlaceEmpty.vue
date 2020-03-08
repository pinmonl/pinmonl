<template>
  <div :class="$style.empty">
    <slot></slot>
  </div>
</template>

<style lang="scss" module>
.empty {
  background-color: #fff;
}
</style>
