<template>
  <Article use-breadcrumb current="Support us" :class="$style.supportUs">
    <div :class="$style.content">
      <div :class="$style.bmc">
        <p>If you love the project, consider to</p>
        <div :class="$style.funding"><BuyMeACoffee /></div>
      </div>
      <div :class="$style.membership">
        <p>Or you can support the project by joining Monler memebership.</p>
        <div :class="$style.funding"><Patreon /></div>
      </div>
    </div>
  </Article>
</template>

<script>
import Article from '@/components/layout/Article.vue'
import BuyMeACoffee from '@/components/funding/BuyMeACoffee.vue'
import Patreon from '@/components/funding/Patreon.vue'

export default {
  components: {
    Article,
    BuyMeACoffee,
    Patreon,
  },
  metaInfo: {
    title: 'Support Us',
  },
}
</script>

<style lang="scss" module>
.content {
  @apply leading-normal;

  p {
    @apply mb-1;
  }
}

.bmc {
  @apply mb-8;
}
</style>
