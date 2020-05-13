<template>
  <Article use-breadcrumb current="Sponsors" :class="$style.sponsors">
    <div :class="$style.content">
      <!--
      <div :class="$style.tie">
        <div :class="$style.tieName">Monler</div>
      </div>
      -->
      <div :class="$style.p">Please support us if you love the project.</div>
      <div :class="$style.note">Become a Monler via <Anchor to="https://www.patreon.com/pinmonl" external color>Patreon</Anchor>.</div>
    </div>
  </Article>
</template>

<script>
import Article from '@/components/layout/Article.vue'

export default {
  components: {
    Article,
  },
  metaInfo: {
    title: 'Sponsors',
  },
}
</script>

<style lang="scss" module>
.content {
  @apply leading-normal;
}

.tie {
  @apply mb-4;
}

.tieName {
  @apply font-bold;
  @apply mb-1;
}

.p {
  @apply mb-2;
}

.note {
  @apply italic;
  @apply text-sm;
}
</style>
