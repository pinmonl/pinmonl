<template>
  <div :class="$style.header">
    <div :class="$style.container">
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
  @apply border-b;
}

.container {
  @apply h-header;
  @apply px-4;
  @apply relative;
  @apply z-10;
  @apply flex;
  @apply items-center;
  @apply justify-between;
  @apply mx-auto;

  @screen lg {
    max-width: 1080px;
  }

  @screen xl {
    max-width: 1300px;
  }
}

.hamburger {
  @apply leading-0;
  @apply mr-2;
  @apply absolute;

  @screen md {
    @apply hidden;
  }
}

.search {
  @apply flex-shrink;
}

.title,
.statusBar {
  @apply flex-shrink-0;
  width: 36px;
}

.account {
  @screen sm-down {
    @apply hidden;
  }
}
</style>
