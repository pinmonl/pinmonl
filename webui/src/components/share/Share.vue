<template>
  <div :class="$style.container">
    <slot name="before"></slot>
    <div :class="$style.shareContainer">
      <div :class="$style.name">
        {{ share.name }}
      </div>
      <div :class="$style.tags">
        <Tag v-for="tag in share.mustTags" :key="tag" :tag="tag" :class="$style.tag" />
      </div>
    </div>
    <Anchor v-if="sharingDest" :to="sharingDest" underline :class="$style.sharing">
      <IconLabel name="link" />
    </Anchor>
  </div>
</template>

<script>
import IconLabel from '@/components/icon/IconLabel.vue'
import Tag from '@/components/tag/Tag.vue'

export default {
  components: {
    IconLabel,
    Tag,
  },
  props: {
    share: {
      type: Object,
      required: true,
    },
  },
  computed: {
    sharingDest () {
      const user = this.$store.getters.user
      const share = this.share
      if (!user || !share.name) {
        return
      }
      const userName = this.$store.getters.user.name
      const shareName = share.name

      return {
        name: 'sharing',
        params: {
          user: userName,
          name: shareName,
        },
      }
    },
  },
  methods: {
  },
}
</script>

<style lang="scss" module>
.container {
  @apply p-4;
  @apply relative;
  @apply flex;
  @apply items-center;

  &:hover {
    @apply bg-bg;
  }
}

.shareContainer {
  @apply flex-grow;
}

.name {
  @apply font-bold;
}

.tags {
  @apply -mx-1;
}

.tag {
  @apply m-1;
}

.sharing {
  @apply relative;
  @apply z-10;
  @apply inline-block;
  @apply py-2;
}
</style>
