<template>
  <div :class="$style.header">
    <Hamburger :class="$style.hamburger" @click="openNav" :active="showNav" />
    <slot></slot>
    <Anchor :to="{ name: 'account' }" :class="$style.account" v-if="!noAccount">
      <IconButton name="accountCircle" />
    </Anchor>
  </div>
</template>

<script>
import Hamburger from './Hamburger.vue'
import IconButton from '@/components/form/IconButton.vue'

export default {
  components: {
    Hamburger,
    IconButton,
  },
  props: {
    noAccount: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    showNav () {
      return this.$store.state.showNav
    },
  },
  methods: {
    openNav () {
      this.$store.commit('SET_NAV', true)
    },
  },
}
</script>

<style lang="scss" module>
.header {
  @apply h-header;
  @apply p-4;
  @apply shadow-b-sm;
  @apply relative;
  @apply z-10;
  @apply flex;
  @apply items-center;

  @screen sm-down {
    @apply h-auto;
  }
}

.hamburger {
  @apply leading-0;
  @apply mr-2;

  @screen md {
    @apply hidden;
  }
}

.account {
  @apply leading-0;

  @screen sm-down {
    @apply hidden;
  }
}
</style>
