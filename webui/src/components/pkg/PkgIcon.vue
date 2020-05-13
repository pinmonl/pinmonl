<template>
  <Icon :class="iconClass" :name="providerIcon" />
</template>

<script>
export default {
  props: {
    pkg: {
      type: Object,
      default: null,
    },
    provider: {
      type: String,
      default: null,
    },
    noHover: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    iconClass () {
      return [this.$style.icon, {
        [this.$style.icon_hover]: !this.noHover,
      }]
    },
    providerName () {
      let name = this.provider
      if (!name && this.pkg) {
        name = this.pkg.provider
      }
      return name
    },
    providerIcon () {
      const name = this.providerName
      return name
    },
  },
}
</script>

<style lang="scss" module>
.icon {
  opacity: 0.8;
}

.icon_hover {
  opacity: 0.5;
  transition: opacity 0.4s ease;

  &:hover {
    opacity: 0.8;
  }
}
</style>
