import style from './style.module.scss'
import CopyableStat from '@/components/stat/CopyableStat.vue'
import PkgIcon from '@/components/pkg/PkgIcon.vue'
import Stat from '@/components/stat/Stat.vue'
import StatIcon from '@/components/stat/StatIcon.vue'
import Tooltip from '@/components/app/Tooltip.vue'

export default {
  components: {
    CopyableStat,
    PkgIcon,
    Stat,
    StatIcon,
    Tooltip,
  },
  props: {
    pkg: {
      type: Object,
      required: true,
    },
  },
  computed: {
    latestTag () {
      return this.findStatValue('tag')
    },
    $baseStyle () {
      return style
    },
  },
  methods: {
    findStat (kind) {
      return this.pkg.stats.find((stat) => stat.kind == kind)
    },
    findStatValue (kind) {
      const stat = this.findStat(kind)
      if (stat) {
        return stat.value
      }
      return null
    }
  },
}
