<template>
  <v-list-item
    @mouseenter="hover = true"
    @mouseleave="hover = false"
  >
    <v-list-item-avatar tile>
      <image-file :image-id="pin.imageId"></image-file>
    </v-list-item-avatar>
    <v-list-item-content :class="contentClasses">
      <v-list-item-title>
        <a :class="$style.title" :href="pin.url" target="_blank" @click.stop="">{{ pin.title }}</a>
      </v-list-item-title>
      <v-list-item-subtitle>
        <tag-chips :tags="pin.tags"></tag-chips>
      </v-list-item-subtitle>
    </v-list-item-content>
    <v-list-item-content :class="pkgClasses">
      <pkg-chips :pkgs="pin.pkgs"></pkg-chips>
    </v-list-item-content>
    <v-list-item-action>
      <v-menu>
        <template #activator="{ on }">
          <v-btn icon @click="on.click">
            <v-icon>{{ mdiDotsVertical }}</v-icon>
          </v-btn>
        </template>
        <v-list dense>
          <v-list-item @click="$emit('edit')">
            <v-list-item-title>
              Edit
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-list-item-action>
  </v-list-item>
</template>

<script>
import { mdiPencil, mdiDotsVertical } from '@mdi/js'
import ImageFile from '@/components/files/ImageFile'
import TagChips from '@/components/tag/TagChips'
import PkgChips from '@/components/pkg/PkgChips'

export default {
  name: 'pin-list-item',
  components: {
    ImageFile,
    TagChips,
    PkgChips,
  },
  props: ['pin'],
  data: () => ({
    hover: false,
    mdiPencil,
    mdiDotsVertical,
  }),
  computed: {
    contentClasses () {
      return {
        [this.$style.content]: true,
      }
    },
    pkgClasses () {
      return {
        [this.$style.pkgs]: true,
        [this.$style.pkgsSm]: this.$vuetify.breakpoint.smAndDown,
      }
    }
  },
  methods: {
  },
}
</script>

<style lang="scss" module>
.content {
  flex-grow: 1;
  flex-shrink: 0;
}

.title {
  &:hover {
    text-decoration: underline;
  }
}

.pkgs {
  flex: 0 1 300px;
}

.pkgsSm {
  flex-basis: auto;
}
</style>
