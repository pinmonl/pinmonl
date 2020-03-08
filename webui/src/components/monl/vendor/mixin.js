import Stat from './Stat.vue'
import Vendor from './Vendor.vue'

export default {
  components: {
    Stat,
    Vendor,
  },
  props: {
    pkg: {
      type: Object,
      required: true,
    },
  },
  computed: {
    vendor () {
      return this.pkg.vendor
    },
    vendorUri () {
      return this.pkg.vendorUri
    },
    stats () {
      return this.pkg.stats.reduce((obj, stat) => {
        const { kind } = stat
        return { ...obj, [kind]: stat }
      }, {})
    },
    isReady () {
      return this.pkg.stats.length > 0
    },
  },
}
