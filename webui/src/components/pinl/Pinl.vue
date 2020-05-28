<template>
  <div :class="containerClasses">
    <slot name="before"></slot>
    <Img :id="pinl.imageId" :class="$style.thumbnail" />
    <div :class="$style.content">
      <div :class="$style.title">
        <Anchor :to="pinl.url" external color>
          {{ pinl.title }}
        </Anchor>
      </div>
      <div :class="$style.description">{{ pinl.description }}</div>
      <div :class="$style.pkgs">
        <template v-for="provider in pkgProviders">
          <PkgIcon
            :key="provider"
            :provider="provider"
            :class="$style.pkg"
          />
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import Img from '@/components/media/Img.vue'
import PkgIcon from '@/components/pkg/PkgIcon.vue'

export default {
  components: {
    Img,
    PkgIcon,
  },
  props: {
    pinl: {
      type: Object,
    },
    active: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    containerClasses () {
      return [this.$style.container, {
        [this.$style.container_active]: this.active,
      }]
    },
    tags () {
      return this.pinl.tags
    },
    pkgProviders () {
      return this.pinl.pkgs.reduce((providers, pkg) => {
        if (!providers.includes(pkg.provider)) {
          return [ ...providers, pkg.provider ]
        }
        return providers
      }, [])
    },
  },
}
</script>

<style lang="scss" module>
$image-size: 50px;

.container {
  @apply p-4;
  @apply relative;

  @include hover-highlight;
}

.container_active {
  @include highlight;
}

.thumbnail {
  @apply float-left;
  width: $image-size;
  height: $image-size;
  padding: 2px;
  @apply rounded-lg;
  @apply object-cover;
}

.content {
  margin-left: $image-size;
  min-height: $image-size;
  @apply px-4;
}

.title {
  @apply inline-block;
  @apply font-bold;

  > a {
    @apply relative;
    @apply z-10;
  }
}

.description {
  @apply mb-1;
  @apply text-xs;
}

.pkgs {
  @apply flex;
  @apply flex-wrap;
}

.pkg {
  @apply relative;
  @apply z-10;
}
</style>
