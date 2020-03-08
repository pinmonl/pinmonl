<template>
  <button :class="buttonClass" type="button" v-on="$listeners">
    <div :class="$style['hamburger-box']">
      <div :class="$style['hamburger-inner']"></div>
    </div>
  </button>
</template>

<script>
export default {
  props: {
    active: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    buttonClass () {
      return [
        this.$style.hamburger, 
        this.$style['hamburger--slider'],
        {
          [this.$style['is-active']]: this.active,
        },
      ]
    },
  },
}
</script>

<style lang="scss" module>
$hamburger-padding-x: 0px;
$hamburger-padding-y: 0px;
$hamburger-layer-width: 24px;
$hamburger-layer-height: 3px;
$hamburger-layer-spacing: 5px;
$hamburger-types: (slider);

@import "~hamburgers/_sass/hamburgers/hamburgers";

.hamburger {
}
</style>
