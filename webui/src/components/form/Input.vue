<template>
  <input
    :class="inputClass"
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
    inputClass () {
      return [this.$style.input, {
        [this.$style.input_error]: this.error,
        [this.$style.input_border]: !this.noStyle,
      }]
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

.input_error {
  &,
  &:hover,
  &:focus {
    @apply border-error;
  }
}
</style>
