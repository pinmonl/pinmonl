<template>
  <input
    :class="[$style.input, {
      [$style.input_border]: !noStyle,
    }]"
    v-bind="props"
    v-on="listeners"
  />
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    type: {
      type: String,
      default: 'text',
    },
    value: {
      default: null,
    },
    error: {
      type: Boolean,
      default: false,
    },
    noStyle: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    props () {
      return {
        ...this.$attrs,
        type: this.type,
        value: this.value,
      }
    },
    listeners () {
      return {
        ...this.$listeners,
        input: (event) => {
          this.$emit('input', event.target.value)
        },
      }
    },
  },
  methods: {
    focus () {
      this.$el.focus()
    },
    blur () {
      this.$el.blur()
    },
  },
}
</script>

<style lang="scss" module>
.input {
  @apply text-sm;
}

.input_border {
  @apply border-solid;
  @apply border;
  @apply border-control;
  @apply px-3;
  @apply h-input;
  transition: all .3s;
  @apply rounded;

  &:hover,
  &:focus {
    @apply border-primary;
  }
}
</style>
