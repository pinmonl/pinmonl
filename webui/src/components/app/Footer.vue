<template>
  <div :class="$style.footer">
    <div :class="$style.container">
      <div :class="$style.item"><Anchor :to="{ name: 'sponsors' }">Sponsors</Anchor></div>
      <div :class="$style.item"><Anchor :to="{ name: 'support-us' }">Support us</Anchor></div>
    </div>
  </div>
</template>

<style lang="scss" module>
.footer {
  @apply m-box;
  @apply p-4;
  @apply border;
  @apply bg-container;
  @apply text-sm;
  @apply rounded;

  @screen xl {
    @apply sticky;
    @apply top-box;
  }
}

.container {
  @apply flex;
  @apply flex-col;
}

.item {
  @apply py-1;

  &,
  & a {
    @apply text-text-secondary;
  }
}
</style>
