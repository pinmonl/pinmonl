<template>
  <Animation appear enter-active-class="slideInRight" leave-active-class="slideOutRight">
    <div :class="$style.container">
      <slot name="before"></slot>
      <div :class="$style.body">
        <slot></slot>
      </div>
      <div :class="$style.control">
        <slot name="control">
          <IconButton @click="$emit('close')" v-if="!noClose" name="close" />
          <IconButton @click="$emit('edit')" v-if="!noEdit" name="edit" />
          <IconButton @click="$emit('cancel')" v-if="!noCancel" name="cancel" />
          <IconButton @click="$emit('save')" v-if="!noSave" name="save" />
        </slot>
        <slot name="other-control"></slot>
      </div>
    </div>
  </Animation>
</template>

<script>
import IconButton from '@/components/form/IconButton.vue'

export default {
  components: {
    IconButton,
  },
  props: {
    noSave: {
      type: Boolean,
      default: false,
    },
    noEdit: {
      type: Boolean,
      default: false,
    },
    noClose: {
      type: Boolean,
      default: false,
    },
    noCancel: {
      type: Boolean,
      default: false,
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply fixed;
  @apply inset-y-0;
  @apply right-0;
  @apply shadow-l;
  @apply z-10;
  @apply flex;
  @apply h-full;
  @apply bg-clear;

  @screen xl {
    width: 600px;
  }

  @screen xxl {
    width: 800px;
  }

  @screen md-down {
    left: theme('width.nav');
  }

  @screen sm-down {
    left: 0;
  }
}

.body {
  @apply flex-1;
  @apply min-w-0;
  @apply relative;
}

.control {
  width: 40px;
  @apply flex-shrink-0;
}
</style>