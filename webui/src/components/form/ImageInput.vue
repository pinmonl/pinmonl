<script>
import ImageFile from '@/components/files/ImageFile'
import { VOverlay, VHover, VFileInput, VImg } from 'vuetify/lib'

export default {
  name: 'image-input',
  props: ['imageId', 'dataUriSrc', 'width', 'height'],
  data: () => ({
    uploadedSrc: null,
  }),
  computed: {
    hasImage () {
      return !!this.imageId || !!this.dataURI
    },
    dataURI () {
      if (this.uploadedSrc) {
        return this.uploadedSrc
      }
      if (this.dataUriSrc) {
        return `data:image/png;base64,${this.dataUriSrc}`
      }
      return null
    },
  },
  methods: {
    genInput (show) {
      const input = this.$createElement(VFileInput, {
        attrs: { hideInput: true, accept: 'image/*', ...this.$attrs },
        on: {
          change: (file) => {
            if (!file) {
              this.$emit('change', file)
              return
            }
            const reader = new FileReader()
            reader.onload = (e) => this.uploadedSrc = e.target.result
            reader.readAsDataURL(file)
            this.$emit('change', file)
          },
        },
      })
      return this.$createElement(VOverlay, {
        attrs: { value: show, absolute: true },
      }, [input])
    },
    genImage (slot) {
      const data = {
        class: 'mb-2',
        attrs: { width: this.width, height: this.height },
      }

      if (this.dataURI) {
        data.attrs.src = this.dataURI
        return this.$createElement(VImg, data, [slot])
      }

      data.attrs.imageId = this.imageId
      return this.$createElement(ImageFile, data, [slot])
    },
  },
  render (h) {
    const defaultSlot = ({ hover }) => {
      const showInput = this.hasImage ? hover : true
      const input = this.genInput(showInput)
      return this.genImage(input)
    }

    return h(VHover, {
      scopedSlots: {
        default: defaultSlot,
      },
    })
  },
}
</script>
