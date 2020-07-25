<template>
  <div>
    <pkg-chip
      v-for="pkg in pkgs"
      :key="pkg.id"
      :pkg="pkg"
    ></pkg-chip>
  </div>
</template>

<script>
import PkgChip from './PkgChip'

export default {
  name: 'pkg-chips',
  components: {
    PkgChip,
  },
  props: ['pkgs'],
}
</script>
