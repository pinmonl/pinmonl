<script>
import * as clipboard from 'clipboard-polyfill'
import Stat from './Stat.vue'
import Tooltip from '@/components/app/Tooltip.vue'

export default {
  props: {
    displayValue: {
      type: String,
    },
    value: {
      type: String,
    },
    tooltip: {
      type: String,
      default: 'Copied',
    },
    delay: {
      type: Number,
      default: 1500,
    },
  },
  data () {
    return {
      showTooltip: false,
      tooltipTimeout: null,
    }
  },
  render (h) {
    const button = () => {
      const displayValue = this.displayValue || this.value
      return h('button', {}, [displayValue])
    }
    const defaultSlot = () => {
      return h(Tooltip, {
        props: {
          text: this.tooltip,
          show: this.showTooltip,
        },
      }, [button])
    }

    return h(Stat, {
      props: { ...this.$props, ...this.$attrs },
      scopedSlots: {
        default: defaultSlot,
      },
      on: {
        click: this.copyToClipboard,
      },
    })
  },
  methods: {
    copyToClipboard () {
      clipboard.writeText(this.value)
      this.show()
    },
    show () {
      this.showTooltip = true
      this.tooltipTimeout && clearTimeout(this.tooltipTimeout)
      this.tooltipTimeout = setTimeout(() => {
        this.hide()
      }, this.delay)
    },
    hide () {
      this.showTooltip = false
      clearTimeout(this.tooltipTimeout)
      this.tooltipTimeout = null
    },
  },
}
</script>
