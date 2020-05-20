<template>
  <button
    :class="buttonClass"
    v-bind="buttonProps"
    v-on="$listeners"
    :disabled="loading"
  >
    <slot name="icon" :iconClass="$style.icon">
      <Icon v-if="loading" name="loading" size="auto" :class="[$style.icon, $style.loading]" />
    </slot>
    <slot></slot>
  </button>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'button',
    },
    noStyle: {
      type: Boolean,
      default: false,
    },
    block: {
      type: Boolean,
      default: false,
    },
    light: {
      type: Boolean,
      default: false,
    },
    danger: {
      type: Boolean,
      default: false,
    },
    warning: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    icon: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    buttonProps () {
      return {
        type: this.type,
      }
    },
    buttonClass () {
      return [{
        [this.$style.button]: !this.noStyle,
        [this.$style.button_block]: this.block,
        [this.$style.button_light]: this.light,
        [this.$style.button_danger]: this.danger,
        [this.$style.button_warning]: this.warning,
        [this.$style.button_loading]: this.loading,
        [this.$style.button_icon]: this.icon,
      }]
    }
  },
}
</script>

<style lang="scss" module>
.button {
  @apply bg-primary;
  @apply text-text-inverted;
  @apply px-4;
  @apply py-1;
  @apply rounded;

  &[disabled] {
    @apply cursor-not-allowed;
  }
}

.button_block {
  @apply block;
}

.button_light {
  @apply bg-btn-light-bg;
  @apply text-btn-light;
}

.button_danger {
  @apply bg-error;
}

.button_warning {
  @apply bg-warning;
}

.button_loading {
  @extend .button_icon;

  &::after {
    content: '';
    @apply absolute;
    @apply inset-0;
    @apply bg-disabled-overlay;
    z-index: 1;
  }
}

.button_icon {
  @apply relative;
  @apply overflow-hidden;
  @apply pl-8;
}

.icon {
  @apply absolute;
  @apply inset-y-0;
  @apply my-auto;
  left: 8px;
  width: 20px;
  height: 20px;
  z-index: 2;
}

.loading {
  animation: spin infinite 1.5s;
}

@keyframes spin {
  0% {
    transform: rotate(0);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
