<template>
  <div :class="containerClass">
    <Animation enter-active-class="slideInLeft" leave-active-class="slideOutLeft">
      <Nav
        v-if="isDesktop ? true : showMobileNav"
        :class="$style.nav"
        :tagSrc="relations.tags"
        :pinlSrc="relations.pinls"
        :parentName="safeParentName"
      />
    </Animation>
    <div :class="$style.body">
      <Header noAccount>
        <h2 :class="$style.title">
          <template v-if="ready">
            {{ owner.login }} / {{ share.name }}
          </template>
        </h2>
        <Search
          :tagSrc="relations.tags"
          v-model="search"
        />
      </Header>
      <template v-if="ready && showInfo && share.description">
        <div :class="$style.description">
          {{ share.description }}
        </div>
        <Divider />
      </template>
      <Pinls :pinlSrc="pinls" :replace="!!pinlId" />
      <RightPanel
        v-if="showPanel"
        @close="handlePanelClose"
        noSave
        noCancel
        noEdit
      >
        <PinlDetail
          :class="$style.pinlDetail"
          :user="user"
          :name="name"
          :id="pinlId"
        />
      </RightPanel>
    </div>
  </div>
</template>

<script>
import { formatRepeatParam } from '@/pkgs/utils'
import Header from '@/components/app/Header.vue'
import Nav from '@/components/sharing/Nav.vue'
import PinlDetail from '@/components/sharing/PinlDetail.vue'
import Pinls from '@/components/sharing/Pinls.vue'
import RightPanel from '@/components/modal/RightPanel.vue'
import Search from '@/components/app/Search.vue'
import { mediaQueries } from '@/theme/variables'

const mdQuery = mediaQueries.get('md')

export default {
  components: {
    Header,
    Nav,
    PinlDetail,
    Pinls,
    RightPanel,
    Search,
  },
  props: {
    user: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    parentName: {
      type: [String, Array],
    },
    pinlId: {
      type: String,
      default: null,
    },
  },
  data () {
    return {
      ready: false,
      loading: false,
      storedParentName: [],
      storedSearch: '',

      share: null,
      relations: {
        pinls: [],
        tags: [],
      },
    }
  },
  created () {
    this.init()
    this.$store.dispatch('mediaQueries/match', mdQuery)
  },
  computed: {
    isDesktop () {
      return this.$store.state.mediaQueries.matches[mdQuery]
    },
    showMobileNav () {
      return this.$store.state.showNav
    },
    selected () {
      return this.safeParentName.length > 0
    },
    owner () {
      return (this.share || {}).owner || {}
    },
    safeParentName () {
      return (this.routeIsRoot || this.routeIsTag)
        ? formatRepeatParam(this.parentName)
        : this.storedParentName
    },
    showInfo () {
      if (this.routeIsRoot) {
        return !this.search
      }
      if (this.routeIsTag) {
        return false
      }
      return this.storedParentName.length == 0
    },
    routeIsRoot () {
      return this.$route.name == 'sharing'
    },
    routeIsTag () {
      return this.$route.name == 'sharing.tag'
    },
    pinls () {
      const { input, tags } = this.$store.getters['pinl/parseSearch'](this.search)
      const tagNames = this.$store.getters['tag/mapName'](tags)
      let pinls = this.relations.pinls
      if (this.safeParentName.length > 0) {
        pinls = this.$store.getters['pinl/getByTag'](pinls, this.safeParentName.slice(-1))
      }
      pinls = this.$store.getters['pinl/getByTag'](pinls, tagNames)
      pinls = this.$store.getters['pinl/searchByTitle'](pinls, input)
      return pinls
    },
    showPanel () {
      return !!this.pinlId
    },
    search: {
      get () {
        if (this.routeIsRoot || this.routeIsTag) {
          return this.$route.query.q
        }
        return this.storedSearch
      },
      set (q) {
        if (q == this.search) {
          return
        }
        this.storedSearch = q
        this.$router.replace({ query: {q} })
      },
    },
    containerClass () {
      return [this.$style.container, {
        [this.$style.container_desktopNav]: this.isDesktop,
        [this.$style.container_mobileNav]: !this.isDesktop && this.showMobileNav,
      }]
    },
  },
  methods: {
    async init () {
      this.share = await this.find(this.user, this.name)
      this.relations.pinls = await this.listPinls(this.user, this.name)
      this.relations.tags = await this.listTags(this.user, this.name)
      this.ready = true
    },
    async find (user, name) {
      return await this.$store.dispatch('sharing/find', { user, name })
    },
    async listPinls (user, name) {
      return await this.$store.dispatch('sharing/listPinls', { user, name })
    },
    async listTags (user, name) {
      return await this.$store.dispatch('sharing/listTags', { user, name })
    },
    handleSelectTag () {
    },
    handlePanelClose () {
      this.$router.push({ name: 'sharing', params: {user: this.user, name: this.name} })
    },
  },
  watch: {
    parentName (newValue, oldValue) {
      if (this.pinlId) {
        this.storedParentName = formatRepeatParam(oldValue)
      } else {
        this.storedParentName = []
      }
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply w-full;
  @apply h-full;
  @apply overflow-y-auto;
}

.container_desktopNav {
  .body {
    margin-left: theme('width.nav');
  }
}

.body {
  @apply h-full;
  
  @screen sm-down {
    @apply ml-0;
  }
}

.nav {
  @apply fixed;
  @apply inset-y-0;
  @apply left-0;
  @apply w-nav;
  @apply z-100;
}

.header {
  @apply p-4;
  min-height: theme('height.header');
  @apply flex;
  @apply items-center;
  @apply shadow-b-sm;
  @apply mb-1;
}

.title {
  @apply flex-grow;
  @apply flex-shrink-0;
}

.description {
  @apply text-sm;
  margin-top: 2px;
  @apply p-4;
}

.pinlDetail {
  @apply p-4;
}
</style>
