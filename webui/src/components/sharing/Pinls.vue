<template>
  <div :class="$style.container">
    <template v-for="pinl in pinls">
      <div :key="pinl.id">
        <Pinl :pinl="pinl">
          <template #before>
            <Anchor :key="pinl.id" :to="dest(pinl)" :replace="replace" inset />
          </template>
        </Pinl>
        <Divider />
      </div>
    </template>
  </div>
</template>

<script>
import Pinl from '@/components/pinl/Pinl.vue'

export default {
  components: {
    Pinl,
  },
  props: {
    pinlSrc: {
      type: Array,
      required: true,
    },
    replace: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    pinls () {
      return this.pinlSrc
    },
  },
  methods: {
    dest (pinl) {
      const { user, name } = this.$route.params
      return {
        name: 'sharing.pinl',
        params: {
          user,
          name,
          pinlId: pinl.id,
        },
      }
    },
  },
}
</script>

<style lang="scss" module>
.container {
}
</style>
