<script>
export default {
  render (h) {
    return h("div", {
      class: [this.$style.box],
    }, [this.$slots.default])
  },
}
</script>

<style lang="scss" module>
.box {
  @apply bg-container;
  @apply shadow-box;
  @apply m-box;
}
</style>
