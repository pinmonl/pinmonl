<template>
  <img :src="imgURL" />
</template>

<script>
export default {
  props: {
    id: {
      type: String,
    },
  },
  computed: {
    imgURL () {
      if (!this.id) {
        return
      }
      return `/api/image/${this.id}`
    },
  },
}
</script>
