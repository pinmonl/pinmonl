<template>
  <div :class="containerClass" @mouseenter="handleShow" @mouseleave="handleHide">
    <slot></slot>
    <div :class="tooltipClass" ref="tooltip">
      <slot name="tooltip">
        <span v-text="text" />
      </slot>
    </div>
  </div>
</template>

<script>
import { createPopper } from '@popperjs/core'
import isBoolean from 'lodash/isBoolean'

export default {
  props: {
    text: {
      type: String,
    },
    show: {
      type: Boolean,
    },
  },
  data () {
    return {
      popper: null,
      localVisible: false,
    }
  },
  mounted () {
    this.init()
  },
  beforeDestroy () {
    this.popper && this.popper.destroy()
    this.popper = null
  },
  computed: {
    containerClass () {
      return [this.$style.container]
    },
    tooltipClass () {
      return [this.$style.tooltip, {
        [this.$style.tooltip_visible]: this.visible,
      }]
    },
    visible: {
      get () {
        if (isBoolean(this.show)) {
          return this.show
        }
        return this.localVisible
      },
      set (value) {
        if (isBoolean(this.show)) {
          this.$emit('update:show', value)
          return
        }
        this.localVisible = value
      },
    },
  },
  methods: {
    init () {
      this.popper = createPopper(this.$el, this.$refs.tooltip, {
        placement: 'top',
      })
    },
    handleShow () {
      this.visible = true
    },
    handleHide () {
      this.visible = false
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply relative;
}

.tooltip {
  @apply invisible;

  background: #333;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
}

.tooltip_visible {
  @apply visible;
}
</style>
