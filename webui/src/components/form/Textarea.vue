<template>
  <textarea
    :class="[$style.textarea, {
      [$style.textarea_border]: !noStyle,
    }]"
    v-bind="props"
    v-on="listeners"
  />
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    value: {
      type: String,
    },
    noStyle: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      height: false,
    }
  },
  mounted () {
    this.updateHeight()
  },
  computed: {
    props () {
      return {
        rows: 1,
        ...this.$attrs,
        style: this.style,
        value: this.value,
      }
    },
    listeners () {
      return {
        ...this.$listeners,
        input: this.handleInput,
      }
    },
    style () {
      const height = this.height ? this.height + 'px' : 'auto'
      return {
        height,
      }
    },
  },
  methods: {
    handleInput (event) {
      this.$emit('input', event.target.value)
    },
    updateHeight () {
      this.height = false
      this.$nextTick(() => {
        this.height = this.$el.scrollHeight
      })
    },
    focus () {
      this.$el.focus()
    },
    blur () {
      this.$el.blur()
    },
  },
  watch: {
    value (newValue, oldValue) {
      if (newValue != oldValue) {
        this.updateHeight()
      }
    },
  },
}
</script>

<style lang="scss" module>
.textarea {
  @apply text-sm;
  @apply resize-none;
  @apply overflow-hidden;
  @apply block;
}

.textarea_border {
  @apply border-solid;
  @apply border;
  @apply border-control;
  @apply px-3;
  @apply py-2;
  min-height: theme('height.input');
  transition: all .3s;
  @apply rounded;

  &:hover,
  &:focus {
    @apply border-primary;
  }
}
</style>
