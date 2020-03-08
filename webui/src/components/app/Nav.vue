<template>
  <div :class="$style.nav">
    <Hamburger @click="closeNav" :active="showNav" :class="$style.hamburger" />
    <div :class="$style.header">
      <Brand to="/" />
    </div>
    <ul :class="$style.navItems">
      <li>
        <Anchor underline to="/bookmark">Bookmark</Anchor>
      </li>
      <li>
        <Anchor underline to="/tag">Tag</Anchor>
      </li>
      <li>
        <Anchor underline to="/share">Share</Anchor>
      </li>
    </ul>
    <div :class="$style.footer">
      <Anchor :to="{ name: 'account' }" :class="$style.accountBtn">
        <IconButton name="accountCircle" block />
      </Anchor>
    </div>
  </div>
</template>

<script>
import Brand from './Brand.vue'
import Hamburger from './Hamburger.vue'
import IconButton from '@/components/form/IconButton.vue'

export default {
  components: {
    Brand,
    Hamburger,
    IconButton,
  },
  computed: {
    showNav () {
      return this.$store.state.showNav
    },
  },
  methods: {
    closeNav () {
      this.$store.commit('SET_NAV', false)
    },
  },
}
</script>

<style lang="scss" module>
.nav {
  @apply flex;
  @apply flex-col;
  @apply shadow-r-sm;
  @apply bg-white;
}

.hamburger {
  @apply absolute;
  top: 1.8rem;
  right: 1rem;

  @screen md {
    @apply hidden;
  }
}

.header {
  @apply h-header;
  @apply p-4;
  @apply flex;
  @apply items-center;
}

.navItems,
.footer {
  @apply list-none;
  @apply m-0;
  @apply p-0;
}

.navItems {
  @apply flex-grow;
  @apply text-sm;

  > li > a {
    @apply block;
    @apply px-6;
    @apply py-2;
  }
}

.footer {
  @apply px-4;
  @apply my-2;
}

.accountBtn {
  @screen md {
    @apply hidden;
  }
}
</style>
