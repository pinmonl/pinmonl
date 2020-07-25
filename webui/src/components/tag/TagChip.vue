<template>
  <span
    :class="tagClasses"
    small 
    label 
    v-bind="$attrs"
  >
    /{{ tagName }}
  </span>
</template>

<script>
export default {
  name: 'tag-chip',
  props: {
    tag: {
      type: [String, Object],
      required: true,
    },
    color: {
      type: String,
      default: 'light',
    },
  },
  computed: {
    tagName () {
      if (typeof this.tag == 'string') {
        return this.tag
      }
      return (this.tag || {}).name
    },
    tagClasses () {
      return {
        [this.$style.tag]: true,
        [this.$style.tagDark]: this.$vuetify.theme.dark,
        [this.$style.fontLight]: this.color == 'light',
      }
    },
  },
}
</script>

<style lang="scss" module>
.tag {
  padding: 0px 4px;
  font-size: 12px;
  box-sizing: border-box;
  display: inline-block;
  background-color: #f0f0f0;
  margin: 2px 2px;
  white-space: nowrap;
  color: rgba(0, 0, 0, 0.8);
  line-height: 1.5;
}

.fontLight {
  color: rgba(0, 0, 0, 0.7);
}

.tagDark {
  background-color: #424242;
  color: rgba(255, 255, 255, 0.8);
}
</style>

