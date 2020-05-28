<template>
  <div :class="$style.container">
    <template v-if="loading">
      <Placeholder>
        <div :class="[$style.phRow, $style.image]"></div>
        <div :class="[$style.phRow, $style.phLine, $style.phTitle]"></div>
        <div :class="[$style.phRow, $style.phLine]"></div>
        <div :class="[$style.phRow, $style.phLine]"></div>
        <div :class="[$style.phRow, $style.phLine, $style.phLineShort]"></div>
      </Placeholder>
    </template>

    <template v-else-if="editable">
      <InputGroup>
        <Label>Title</Label>
        <Input v-model="model.title" />
      </InputGroup>
      <InputGroup>
        <Label>URL</Label>
        <Input v-model="model.url" :error="$v.model.url.$error" ref="url" />
        <template #errors v-if="$v.model.url.$error">
          <p v-if="!$v.model.url.required">URL cannot be empty.</p>
          <p v-if="!$v.model.url.url">Invalid URL.</p>
        </template>
      </InputGroup>
      <InputGroup>
        <Label>Description</Label>
        <Textarea v-model="model.description" />
      </InputGroup>
      <InputGroup>
        <Label>Tags</Label>
        <TagInput v-model="tags" :class="$style.tags" />
      </InputGroup>
    </template>

    <template v-else>
      <Img :class="$style.image" :id="model.imageId" v-if="model.imageId" />
      <div :class="$style.title">
        <Anchor :to="model.url" external color v-text="model.title" />
      </div>
      <div :class="$style.description" v-text="model.description" />
      <TagInput :class="$style.tags" :value="tags" disabled noStyle />
      <div :class="$style.pkgs">
        <template v-for="(pkgs, provider) in groupedPkgs">
          <div :key="provider" :class="$style.pkgProvider">
            <PkgIcon :provider="provider" :class="$style.pkgProviderIcon" />
            <div :class="$style.pkgs">
              <Pkg v-for="pkg in pkgs" :key="pkg.id" :pkg="pkg" :class="$style.pkg" />
            </div>
          </div>
        </template>
      </div>
    </template>

    <slot
      name="controls"
      :error="$v.$error"
      :submit="handleSubmit"
      :cancel="handleCancel"
    />
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import modelMixin from '@/mixins/model'
import placeholderMixin from '@/mixins/placeholder'
import keybindingMixin from '@/mixins/keybinding'
import Img from '@/components/media/Img.vue'
import Pkg from '@/components/pkg/Pkg.vue'
import PkgIcon from '@/components/pkg/PkgIcon.vue'
import TagInput from '@/components/tag/TagInput.vue'
import { validationMixin } from 'vuelidate'
import { required, url } from 'vuelidate/lib/validators'

export default {
  mixins: [formMixin, modelMixin({ prop: 'pinl' }), placeholderMixin, validationMixin, keybindingMixin()],
  components: {
    Img,
    Pkg,
    PkgIcon,
    TagInput,
  },
  props: {
    pinl: {
      type: Object,
      default: null,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  mounted () {
    this.editable && this.autoFocusURL()
  },
  computed: {
    model () {
      const gt = this.getter
      const st = this.setter
      return {
        $data: this.pinl,
        get url () { return gt('url')},
        set url (v) { st('url', v)},
        get title () { return gt('title') },
        set title (v) { st('title', v) },
        get description () { return gt('description')},
        set description (v) { st('description', v)},
        get readme () { return gt('readme')},
        set readme (v) { st('readme', v)},
        get tags () { return gt('tags') },
        set tags (v) { st('tags', v) },
        get pkgs () { return gt('pkgs') },
        set pkgs (v) { st('pkgs', v) },
        get imageId () { return gt('imageId')},
      }
    },
    tags: {
      get () {
        const tags = this.$store.getters['tag/tags']
        return this.$store.getters['tag/getByName'](tags, this.model.tags)
      },
      set (tags) {
        this.model.tags = this.$store.getters['tag/mapName'](tags)
      },
    },
    groupedPkgs () {
      return this.model.pkgs.reduce((pkgs, pkg) => {
        const { provider } = pkg
        if (!pkgs[provider]) {
          pkgs[provider] = []
        }
        pkgs[provider].push(pkg)
        return pkgs
      }, {})
    },
  },
  methods: {
    handleSubmit () {
      this.$v.$touch()
      if (this.$v.$error) {
        return
      }
      this.syncModel()
    },
    handleCancel () {
      this.$v.$reset()
      this.revertModel()
      this.$emit('update:editable', false)
      this.$emit('cancel')
    },
    handleKeyPress (e) {
      if (this.editable) {
        if (e.key == 's' && e.ctrlKey) {
          e.preventDefault()
          this.handleSubmit()
          return false
        }

        return
      }

      if (e.key == 'e') {
        this.$emit('update:editable', true)
        e.preventDefault()
        return false
      }
      if (e.key == 'o') {
        this.$store.dispatch('pinl/openLink', { pinl: this.pinl })
        return
      }
    },
    autoFocusURL () {
      this.$refs.url.focus()
    },
  },
  watch: {
    'editable' (val) {
      if (val) {
        this.$nextTick(() => {
          this.autoFocusURL()
        })
      }
    },
  },
  validations () {
    if (!this.editable) {
      return {}
    }
    return {
      model: {
        title: {},
        url: { required, url },
        description: {},
        readme: {},
        tags: {},
      },
    }
  },
}
</script>

<style lang="scss" module>
.container {
  @apply flex;
  @apply flex-col;
  @apply text-sm;
}

.image {
  width: 80px;
  height: 80px;
  @apply rounded-lg;
  @apply overflow-hidden;
  @apply object-cover;
}

.title {
  @apply font-bold;
  @apply mt-1;
  @apply text-base;
}

.image,
.title,
.url,
.description,
.tags,
.pkgs,
.divider,
.readme {
  @apply mb-2;
}

.phRow {
  @apply mb-2;
}

.phLine {
  height: 1rem;
}

.phTitle {
  width: 50%;
}

.phLineShort {
  width: 30%;
}

.pkgProvider {
  @apply flex;
}

.pkgProviderIcon {
  @apply flex-shrink-0;
}

.pkgs {
  @apply flex-grow;
}

.pkg {
  @apply block;
  line-height: 24px;
}
</style>
