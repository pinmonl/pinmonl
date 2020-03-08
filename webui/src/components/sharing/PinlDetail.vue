<template>
  <PinlDetail
    :loading="loading"
    :pinl="pinl"
  />
</template>

<script>
import PinlDetail from '@/components/pinl/PinlDetail.vue'

export default {
  components: {
    PinlDetail,
  },
  props: {
    id: {
      type: String,
    },
    user: {
      type: String,
    },
    name: {
      type: String,
    },
  },
  data () {
    return {
      loading: false,
      pinl: null,
    }
  },
  created () {
    this.initPinl()
  },
  computed: {
  },
  methods: {
    async initPinl () {
      this.loading = true
      let pinl = null
      if (this.id) {
        pinl = await this.find(this.user, this.name, this.id)
      }

      this.pinl = pinl
      this.loading = false
    },
    async find (user, name, id) {
      return await this.$store.dispatch('sharing/findPinl', { user, name, id })
    },
  },
  watch: {
    id () {
      this.initPinl()
    },
  },
}
</script>
