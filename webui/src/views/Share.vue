<template>
  <div :class="$style.container">
    <Header :class="$style.header">
      <h2 :class="$style.title">Share</h2>
      <Anchor :to="{ name: 'share.new' }" :class="$style.addBtn">
        <IconButton name="add" block />
      </Anchor>
    </Header>
    <div :class="$style.list">
      <template v-for="share in shares">
        <Anchor :to="{ name: 'share', params: {id: share.id} }" :key="share.id" :replace="showPanel">
          <Share :share="share" />
          <Divider />
        </Anchor>
      </template>
    </div>
    <RightPanel
      v-if="showPanel"
      @close="handlePanelClose"
      @save="handlePanelSave"
      noCancel
      noEdit
    >
      <ShareDetail
        v-model="model"
        :loading="loading"
        :editable="!loading"
      />
    </RightPanel>
  </div>
</template>

<script>
import Header from '@/components/app/Header.vue'
import IconButton from '@/components/form/IconButton.vue'
import RightPanel from '@/components/modal/RightPanel.vue'
import Share from '@/components/share/Share.vue'
import ShareDetail from '@/components/share/ShareDetail.vue'

export default {
  components: {
    Header,
    IconButton,
    RightPanel,
    Share,
    ShareDetail,
  },
  props: {
    id: {
      type: String,
      default: null,
    },
    new: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      loading: false,
      editing: false,

      model: null,
      original: null,
    }
  },
  created () {
    this.initModel()
  },
  computed: {
    shares () {
      return this.$store.getters['share/shares']
    },
    hasId () {
      return !!this.id
    },
    showPanel () {
      return this.new || this.hasId
    },
  },
  methods: {
    async initModel () {
      this.loading = true
      if (this.hasId) {
        this.original = await this.find(this.id)
      } else if (this.new) {
        this.original = this.$store.getters['share/new']()
      }
      this.model = this.original
      this.loading = false
    },
    async find(id) {
      return await this.$store.dispatch('share/find', { id })
    },
    async handlePanelSave () {
      const model = await this.save(this.model)
      this.model = model
    },
    async save(model) {
      if (model.id) {
        return await this.$store.dispatch('share/update', model)
      } else {
        const share = await this.$store.dispatch('share/create', model)
        this.$router.replace({ name: 'share', params: {id: share.id} })
        return share
      }
    },
    handlePanelClose () {
      this.$router.push({ name: 'share' })
    },
  },
  watch: {
    id () {
      this.initModel()
    },
    new () {
      this.initModel()
    },
  }
}
</script>

<style lang="scss" module>
.container {
  @apply relative;
  @apply h-full;
  @apply w-full;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}

.header {
  @apply flex;
  @apply flex-wrap;
}

.title {
  @apply flex-grow;
}

.addBtn {
}
</style>
