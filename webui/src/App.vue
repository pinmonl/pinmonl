<template>
  <div :class="appClass">
    <IconInstall />
    <Animation enter-active-class="slideInLeft" leave-active-class="slideOutLeft">
      <Nav :class="$style.nav" v-if="isDesktop ? showNav : showMobileNav" />
    </Animation>
    <div :class="$style.body">
      <router-view/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/app/Nav.vue'
import { mapGetters, mapActions, mapMutations } from 'vuex'
import IconInstall from '@/components/icon/IconInstall.vue'
import { mediaQueries } from '@/theme/variables'

const mdQuery = mediaQueries.get('md')

export default {
  components: {
    Nav,
    IconInstall,
  },
  async created () {
    await this.init()
    this.$store.dispatch('mediaQueries/match', mdQuery)
  },
  computed: {
    isDesktop () {
      return this.$store.state.mediaQueries.matches[mdQuery]
    },
    showNav () {
      return this.$route.matched.some(m => m.meta.showNav)
    },
    showMobileNav () {
      return this.$store.state.showNav
    },
    appClass () {
      return [this.$style.app, {
        [this.$style.app_desktopNav]: this.isDesktop && this.showNav,
        [this.$style.app_mobileNav]: !this.isDesktop && this.showMobileNav,
      }]
    },
    ...mapGetters([
      'authed',
    ]),
  },
  methods: {
    async init () {
      try {
        await this.getMe()
      } catch (e) {
        //
      } finally {
        this.setAppReady(true)
      }
    },
    async fetchUserData () {
      await Promise.all([
        this.$store.dispatch('pinl/fetchAll'),
        this.$store.dispatch('tag/fetchAll'),
        this.$store.dispatch('share/fetchAll')
      ])
    },
    ...mapMutations({
      setAppReady: 'SET_READY',
    }),
    ...mapActions([
      'getMe',
    ]),
  },
  watch: {
    authed (authed) {
      if (authed) {
        this.fetchUserData()
      }
    },
  },
}
</script>

<style lang="scss" module>
@tailwind base;
@import './theme/base.scss';

.app {
  @apply h-full;
  @apply w-full;
  @apply text-text-primary;
  @apply overflow-hidden;
}

.app_desktopNav {
  .body {
    margin-left: theme('width.nav');
  }
}

.nav {
  @apply fixed;
  @apply inset-y-0;
  @apply left-0;
  @apply z-100;
  width: theme('width.nav');
}

.body {
  @apply h-full;

  @screen sm-down {
    @apply ml-0;
  }
}
</style>
