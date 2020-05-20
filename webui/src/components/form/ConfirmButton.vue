<script>
import Button from './Button.vue'

export default {
  inheritAttrs: false,
  props: {
    asking: {
      type: Boolean,
      default: false,
    },
    timeout: {
      type: Number,
      default: 2000,
    },
  },
  data () {
    return {
      timer: null,
    }
  },
  beforeDestroy () {
    this.clearTimer()
    this.cancel()
  },
  render (h) {
    return h(Button, {
      props: {
        warning: this.asking,
        icon: this.asking,
        ...this.$attrs,
      },
      on: {
        click: () => {
          if (this.asking) {
            this.confirm()
          } else {
            this.touch()
          }
        },
      },
      scopedSlots: {
        icon: (slotProps) => {
          if (!this.asking) {
            return
          }
          return h('Icon', {
            class: slotProps.iconClass,
            props: { name: 'exclamation', size: 'auto' },
          })
        },
      },
    }, [
      !this.asking && this.$slots.default,
      this.asking && this.renderConfirmation(h),
    ])
  },
  methods: {
    renderConfirmation () {
      return 'Click to confirm'
    },
    touch () {
      this.autoCancel()
      this.$emit('update:asking', true)
    },
    confirm () {
      this.$emit('update:asking', false)
      this.$emit('click')
    },
    cancel () {
      this.$emit('update:asking', false)
    },
    autoCancel () {
      this.clearTimer()
      this.timer = setTimeout(() => {
        this.cancel()
      }, this.timeout)
    },
    clearTimer () {
      this.timer && clearTimeout(this.timer)
    },
  },
  watch: {
    'asking' (newValue, oldValue) {
      if (newValue == oldValue) {
        return
      }
      if (newValue) {
        this.autoCancel()
      } else {
        this.clearTimer()
      }
    },
  },
}
</script>

<style lang="scss" module>
.confirmation {
  @apply relative;
  @apply pl-4;
}

.confirmationIcon {
  @apply absolute;
  @apply inset-y-0;
  @apply my-auto;
  left: 8px;
  width: 18px;
  height: 18px;
}
</style>
