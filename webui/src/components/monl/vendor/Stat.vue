<template>
  <div :class="$style.stat">
    <div :class="$style.icon">
      <Icon :name="icon" size="fill" />
    </div>
    <div :class="$style.value">{{ formattedValue }}</div>
  </div>
</template>

<script>
import { numberFormat } from '@/pkgs/format'

export default {
  props: {
    icon: { 
      type: String,
    },
    value: {
      type: [String, Number],
    },
    number: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    formattedValue () {
      if (this.number) {
        return numberFormat(this.value)
      }
      return this.value
    },
  },
}
</script>

<style lang="scss" module>
.stat {
  @apply flex;
  @apply items-center;
  @apply text-xs;
  @apply mr-3;
}

.icon {
  width: 1rem;
  height: 1rem;
  @apply mr-1;
}

.value {
}
</style>
