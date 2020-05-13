<template>
  <Anchor :class="$style.button" to="https://www.patreon.com/pinmonl" external>
    <img :src="imageSrc" />
  </Anchor>
</template>

<script>
import patronButton from '@/assets/become_a_patron_button.png'

export default {
  computed: {
    imageSrc () {
      return patronButton
    },
  },
}
</script>

<style lang="scss" module>
.button {
  @apply inline-block;
  @apply overflow-hidden;
  @apply rounded;
}
</style>
