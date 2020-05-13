<template>
  <div :class="$style.modal">
    <div :class="$style.container">
      <Backdrop @click="handleBackdrop" :class="$style.backdrop" absolute />
      <div :class="$style.content">
        <button :class="$style.close" @click="handleClose">
          <Icon name="close" />
        </button>
        <slot name="header" :headerClass="$style.header"></slot>
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    disableKeys: {
      type: Boolean,
      default: false,
    },
  },
  mounted () {
    document.addEventListener("keyup", this.handleKeyPress)
  },
  beforeDestroy () {
    document.removeEventListener("keyup", this.handleKeyPress)
  },
  methods: {
    handleBackdrop () {
      this.$emit("backdrop")
    },
    handleClose () {
      this.$emit("close")
    },
    handleKeyPress (e) {
      if (this.disableKeys) {
        return
      }

      if (e.key == "Escape") {
        this.handleClose()
        return
      }
    },
  },
}
</script>

<style lang="scss" module>
.modal {
  @apply fixed;
  @apply inset-0;
  @apply z-500;
  @apply overflow-y-auto;
  @apply scrolling-touch;
}

.container {
  @apply w-full;
  @apply min-h-full;
  @apply relative;
  @apply flex;
  @apply flex-wrap;
  @apply items-center;
  @apply justify-center;
}

.backdrop {
  @apply bg-backdrop;
}

.close {
  @apply absolute;
  top: 0;
  right: 0;
  @apply p-1;
  @apply z-10;
}

.content {
  @apply relative;
  @apply bg-container;
  @apply m-box;
  @apply rounded;
  @apply p-4;

  max-width: 600px;
  width: 600px;
  @apply shadow-modal;
}

.header {
  @apply uppercase;
  @apply font-bold;
  @apply mb-4;
  @apply mx-1;
  @apply text-xs;
}
</style>
