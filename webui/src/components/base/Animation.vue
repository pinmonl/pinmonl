<script>
export default {
  inheritAttrs: false,
  props: {
    enterActiveClass: {
      type: String,
    },
    enterClass: {
      type: String,
    },
    leaveActiveClass: {
      type: String,
    },
    leaveClass: {
      type: String,
    },
    type: {
      default: 'animation',
    },
    animationDuration: {
      default: 300,
    },
    group: {
      type: Boolean,
      defualt: false,
    },
  },
  render (h) {
    const c = this.group ? 'transition-group' : 'transition'
    return h(c, {
      props: this.props,
      on: this.listeners,
    }, this.$slots.default)
  },
  computed: {
    props () {
      const props = { ...this.$attrs }
      for (const [k, v] of Object.entries(this.$props)) {
        if (typeof v == 'undefined') {
          continue
        }
        if (/class$/i.test(k)) {
          props[k] = this.$style.animate + ' ' + this.$style[v]
        } else {
          props[k] = v
        }
      }
      return props
    },
    hasEnter () {
      const re = /enter.*class$/i
      return Object.keys(this.$props).filter(k => re.test(k)).length > 0
    },
    hasLeave () {
      const re = /leave.*class$/i
      return Object.keys(this.$props).filter(k => re.test(k)).length > 0
    },
    listeners () {
      const events = {}
      const ad = this.animationDuration
      let adEnter, adLeave
      adEnter = adLeave = ad
      if (typeof ad == 'object') {
        adEnter = ad.enter
        adLeave = ad.Leave
      }

      if (this.hasEnter) {
        events.beforeEnter = (el) => {
          el.style.animationDuration = `${adEnter}ms`
        }
        events.afterEnter = events.enterCancelled = (el) => {
          el.style.animationDuration = ''
        }
      }
      if (this.hasLeave) {
        events.beforeLeave = (el) => {
          el.style.animationDuration = `${adLeave}ms`
        }
        events.afterLeave = events.leaveCancelled = (el) => {
          el.style.animationDuration = ''
        }
      }
      return events
    },
  },
}
</script>

<style lang="scss" module>
@import '~animate.css/source/sliding_entrances/slideInLeft';
@import '~animate.css/source/sliding_entrances/slideInRight';
@import '~animate.css/source/sliding_exits/slideOutLeft';
@import '~animate.css/source/sliding_exits/slideOutRight';

.animate {
  animation-fill-mode: both;
}
</style>
