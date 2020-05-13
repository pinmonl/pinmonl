<template>
  <div :class="$style.threeColumn">
    <div :class="$style.container">
      <div :class="$style.left">
        <slot name="left"></slot>
      </div>
      <div :class="$style.centered">
        <slot></slot>
      </div>
      <div :class="$style.right">
        <slot name="right"></slot>
      </div>
    </div>
  </div>
</template>

<script>
import Footer from './Footer.vue'

export default {
  components: {
    Footer,
  },
}
</script>

<style lang="scss" module>
.threeColumn {
  @apply h-full;
  @apply w-full;
  @apply overflow-auto;
  @apply scrolling-touch;
}

.container {
  @apply w-full;
  @apply flex;
  @apply mx-auto;

  @screen md-down {
    @apply flex-col;
  }

  @screen lg {
    max-width: 720px + 200px*2; 
  }

  @screen xl {
    max-width: 900px + 240px*2;
  }
}

.left,
.right {
  @apply relative;
  @apply flex-shrink-0;

  @screen lg {
    width: 200px;
  }

  @screen xl {
    width: 240px;
  }
}

.centered {
  @screen lg {
    width: 720px;
  }
  @screen xl {
    width: 900px;
  }
}

.right {
  @screen md-down {
    @apply hidden;
  }

  @screen lg {
    grid-area: right;
  }
}
</style>
