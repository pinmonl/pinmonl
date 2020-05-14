<template>
  <div :class="$style.container">
    <template v-for="(item, i) in items">
      <Anchor
        :key="i"
        :class="$style.anchor"
        :active-class="item.noHighlight ? '' : $style.anchor_active"
        :to="item.to"
        :exact="item.exact"
      >
        <div :class="$style.label">
          <template v-if="item.label">
            {{ item.label }}
          </template>
          <template v-else-if="item.component">
            <component :is="item.component" />
          </template>
        </div>
      </Anchor>
    </template>
    <slot
      name="controls"
      :anchorClass="$style.anchor"
      :anchorActiveClass="$style.anchor_active"
      :labelClass="$style.label"
    />
  </div>
</template>

<script>
const navItems = [
  {
    to: '/bookmark',
    icon: 'bookmark',
    label: 'Bookmark',
  },
  {
    to: '/tag',
    icon: 'tag',
    label: 'Tag',
  },
]

export default {
  computed: {
    items () {
      return [ ...navItems ]
    },
  },
}
</script>

<style lang="scss" module>
$icon-size: 16px;

@mixin anchor-bar-active {
  @screen md-down {
    width: 100%;
    left: 0;
  }

  @screen lg {
    height: 100%;
    top: 0;
  }
}

.container {
  @apply flex;
  @apply text-sm;
  @apply mt-box;
  @apply mx-box;

  @screen lg {
    @apply flex-col;
    @apply items-end;
  }

  @apply relative;
  @apply top-0;
  @screen xl {
    @apply sticky;
    @apply top-box;
  }
}

.anchor {
  @apply relative;
  padding: 8px 12px;

  &:before {
    content: "";
    @apply block;
    @apply absolute;
    background: #00000080;
  }
  &:hover:before {
    @include anchor-bar-active;
  }

  @screen md-down {
    margin: 0 2px;
    &:before {
      transition: width 0.3s ease, left 0.3s ease;
      height: 2px;
      width: 0;
      left: 50%;
      bottom: 0;
    }
  }

  @screen lg {
    @apply text-right;
    margin: 4px 0;
    &:before {
      transition: height 0.3s ease 0s, top 0.3s ease 0s;
      height: 0;
      width: 2px;
      top: 50%;
      right: 0;
    }
  }
}

.anchor_active {
  &:before {
    @apply bg-primary;
    @include anchor-bar-active;
  }
}

.label {
}
</style>
