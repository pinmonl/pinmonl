<template>
  <div :class="$style.header">
    <!-- <Hamburger :class="$style.hamburger" @click="openNav" :active="showNav" /> -->
    <div :class="$style.title">
      <Anchor to="/">
        <LogoSVG />
      </Anchor>
    </div>
    <div :class="$style.centeredContent">
      <slot></slot>
    </div>
    <div :class="$style.statusBar">
      <!--
      <Anchor :to="{ name: 'account' }" :class="$style.account" v-if="!noAccount">
        <IconButton name="accountCircle" block />
      </Anchor>
      -->
    </div>
  </div>
</template>

<script>
import LogoSVG from '@/assets/logo.svg'

export default {
  components: {
    LogoSVG,
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
  @apply px-4;
  @apply relative;
  @apply z-10;
  @apply flex;
  @apply items-center;
  @apply justify-between;
  @apply border-b;
}

.hamburger {
  @apply leading-0;
  @apply mr-2;
  @apply absolute;

  @screen md {
    @apply hidden;
  }
}

.title {
  @apply flex-shrink-0;
  width: 36px;
}

.search {
  @apply flex-shrink;
}

.statusBar {
  @apply flex-shrink;
}

.account {
  @screen sm-down {
    @apply hidden;
  }
}
</style>
