<template>
  <div :class="$style.container">
    <div :class="[$style.icon, iconClass]">
      <Icon :name="name" />
    </div>
    <div :class="[$style.label, labelClass]">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    name: String,
    iconClass: [String, Object, Array],
    labelClass: [String, Object, Array],
  },
}
</script>

<style lang="scss" module>
.container {
  @apply flex;
  @apply items-center;
}

.icon {
  @apply mr-1;
}
</style>
