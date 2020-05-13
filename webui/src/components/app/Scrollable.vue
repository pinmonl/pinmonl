<script>
export default {
  render (h) {
    return h('div', {
      class: this.$style.scrollable,
    }, this.$slots.default)
  },
}
</script>

<style lang="scss" module>
.scrollable {
  @apply w-full;
  @apply h-full;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}
</style>
