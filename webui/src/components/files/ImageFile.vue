<script>
import { VImg } from 'vuetify/lib'
import { baseURL } from '@/utils/constants'

export default {
  name: 'image-file',
  props: {
    imageId: {
      type: String,
    },
  },
  computed: {
    hasImage () {
      return !!this.imageId
    },
    src () {
      if (!this.hasImage) {
        return
      }
      return `${baseURL}/image/${this.imageId}`
    },
  },
  render (h) {
    return h(VImg, {
      attrs: { ...this.$attrs, src: this.src },
      on: { ...this.$listeners },
    }, [this.$slots.default])
  },
}
</script>
