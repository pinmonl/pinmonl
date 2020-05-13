<template>
  <div :class="$style.container" v-on="$listeners">
    <slot name="icon" :classes="{ icon: $style.icon }">
      <div v-if="kind" :class="$style.icon">
        <StatIcon :kind="kind" />
      </div>
    </slot>
    <slot :valueClass="$style.value" :value="formattedValue">
      <div :class="$style.value">{{ formattedValue }}</div>
    </slot>
  </div>
</template>

<script>
import StatIcon from './StatIcon.vue'
import { 
  numberFormat, 
  humanFileSize,
  humanizeNumber,
} from '@/pkgs/format'

export default {
  components: {
    StatIcon,
  },
  props: {
    kind: {
      type: String,
    },
    value: {
      type: String,
    },
    asNumber: {
      type: Boolean,
      default: false,
    },
    asFileSize: {
      type: Boolean,
      default: false,
    },
    asLargeNumber: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    formattedValue () {
      if (this.asNumber) {
        return numberFormat(this.value)
      }
      if (this.asFileSize) {
        return humanFileSize(this.value, true)
      }
      if (this.asLargeNumber) {
        return humanizeNumber(this.value)
      }
      return this.value
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply flex;
  @apply items-center;
}

.icon {
  margin-right: 1px;
}

.value {
}
</style>
