<template>
  <div :class="classes" :style="styles">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'app-container',
  computed: {
    classes () {
      return ['mx-auto']
    },
    styles () {
      return { 'max-width': '900px' }
    },
  },
}
</script>
