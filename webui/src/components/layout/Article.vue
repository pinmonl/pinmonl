<script>
import Box from '@/components/app/Box.vue'
import Breadcrumb from '@/components/app/Breadcrumb.vue'
import Scrollable from '@/components/app/Scrollable.vue'

export default {
  inheritAttrs: false,
  props: {
    useBreadcrumb: {
      type: Boolean,
      default: false,
    },
  },
  render (h) {
    return h(Scrollable, null, [
      h('div', {
        class: this.$style.container,
      }, [
        this.renderBreadcrumb(h),
        h(Box, null, [
          h('div', {
            class: this.$style.content,
          }, this.$slots.default)
        ])
      ])
    ])
  },
  methods: {
    renderBreadcrumb (h) {
      if (!this.useBreadcrumb) {
        return null
      }
      return h(Breadcrumb, {
        props: this.$attrs,
      })
    },
  },
}
</script>

<style lang="scss" module>
.container {
  max-width: 720px;
  @apply mx-auto;
}

.content {
  @apply py-8;
  @apply px-4;
  @apply text-sm;

  @screen md {
    @apply px-8;
    @apply text-base;
  }
}
</style>
