<template>
  <div :class="$style.nav">
    <Hamburger @click="closeNav" :active="showNav" :class="$style.hamburger" />
    <div :class="$style.header">
      <Brand :to="homeRoute" />
    </div>
    <div :class="$style.parents" v-if="parents.length">
      <template v-for="(tag, n) in parents">
        <component
          :key="tag.id"
          :is="isLastParent(tag, n) ? 'div' : 'Anchor'"
          v-bind="isLastParent(tag, n) ? {
          } : {
            color: true,
            to: { name: 'sharing.tag', params: { parentName: parentsAt(n) } }
          }"
        >
          <div :class="$style.tag">
            {{ tag.name }}
          </div>
        </component>
      </template>
    </div>
    <div :class="$style.tags">
      <template v-for="tag in tags">
        <Anchor :to="destOf(tag)" :key="tag.id" color>
          <div :class="$style.tag">
            {{ tag.name }}
          </div>
        </Anchor>
      </template>
    </div>
  </div>
</template>

<script>
import Brand from '@/components/app/Brand.vue'
import Hamburger from '@/components/app/Hamburger.vue'

export default {
  components: {
    Brand,
    Hamburger,
  },
  props: {
    tagSrc: {
      type: Array,
      default: () => ([]),
    },
    parentName: {
      type: Array,
      default: () => ([]),
    },
    pinlSrc: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    parents () {
      return this.$store.getters['tag/getByName'](this.tagSrc, this.parentName)
    },
    parent () {
      const len = this.parents.length
      if (!len) {
        return null
      }
      return this.parents[len - 1]
    },
    parentId () {
      return this.parent ? this.parent.id : ''
    },
    tags () {
      return this.$store.getters['tag/getByParent'](this.tagSrc, this.parentId)
    },
    homeRoute () {
      return {
        name: 'sharing', 
        params: {
          user: this.$route.params.user,
          name: this.$route.params.name,
        },
      }
    },
    showNav () {
      return this.$store.state.showNav
    },
  },
  methods: {
    destOf (tag) {
      const params = { ...this.$route.params, pinlId: null }
      params.parentName = [ ...this.parentName, tag.name ]
      return { name: 'sharing.tag', params }
    },
    parentsAt (n) {
      return this.parentName.slice(0, n + 1)
    },
    isLastParent (tag, n) {
      return (n + 1) == this.parents.length
    },
    closeNav () {
      this.$store.commit('SET_NAV', false)
    },
  },
}
</script>

<style lang="scss" module>
.nav {
  @apply overflow-y-auto;
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

.parents {
  @apply mb-1;
  @apply font-bold;
  @apply mx-2;
}

.tags {
  @apply ml-6;
  border-left: dotted 0.25rem;
  @apply border-anchor;
}

.tag {
  @apply text-sm;
  @apply px-3;
  @apply py-1;
}
</style>
