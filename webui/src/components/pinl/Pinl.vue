<template>
  <div :class="$style.container">
    <slot name="before"></slot>
    <Img :id="pinl.imageId" :class="$style.thumbnail" />
    <div :class="$style.content">
      <div :class="$style.title">
        <Anchor :to="pinl.url" external color>
          {{ pinl.title }}
        </Anchor>
      </div>
      <!-- <div :class="$style.description">{{ pinl.description }}</div> -->
      <div :class="$style.tags">
        <Tag v-for="tag in tags" :key="tag" :tag="tag" :class="$style.tag" />
      </div>
    </div>
  </div>
</template>

<script>
import Img from '@/components/media/Img.vue'
import Tag from '@/components/tag/Tag.vue'

export default {
  components: {
    Img,
    Tag,
  },
  props: {
    pinl: {
      type: Object,
    },
  },
  computed: {
    tags () {
      return this.pinl.tags
    },
  },
}
</script>

<style lang="scss" module>
$image-size: 50px;

.container {
  @apply p-4;
  @apply relative;
  @apply relative;

  &:hover {
    @apply bg-bg;
  }
}

.thumbnail {
  @apply float-left;
  width: $image-size;
  height: $image-size;
  padding: 2px;
  @apply rounded-lg;
}

.content {
  margin-left: $image-size;
  min-height: $image-size;
  @apply px-4;
}

.title {
  @apply inline-block;
  @apply font-bold;
  @apply relative;
  @apply z-10;
}

.description {
  @apply mb-1;
  @apply truncate;
}

.tags {
  @apply -mx-1;
}

.tag {
  @apply m-1;
}
</style>
