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
        <Input v-model="model.url" />
      </InputGroup>
      <InputGroup>
        <Label>Description</Label>
        <Textarea v-model="model.description" />
      </InputGroup>
      <InputGroup>
        <Label>Tags</Label>
        <TagInput v-model="tags" :class="$style.tags" />
      </InputGroup>
      <InputGroup>
        <Label>Readme</Label>
        <Textarea ref="readme" v-model="model.readme" :class="$style.readme" />
      </InputGroup>
    </template>

    <template v-else>
      <div :class="$style.image" v-if="model.imageId">
        <Img :id="model.imageId" />
      </div>
      <div :class="$style.title">
        <Anchor :to="model.url" external color v-text="model.title" />
      </div>
      <div :class="$style.description" v-text="model.description" />
      <TagInput :class="$style.tags" :value="tags" disabled noStyle />
      <div :class="$style.pkgs">
        <Pkg v-for="pkg in model.pkgs" :key="pkg.id" :pkg="pkg" />
      </div>
      <Divider :class="$style.divider" />
      <div v-text="model.readme" :class="$style.readme" />
    </template>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import modelMixin from '@/mixins/model'
import placeholderMixin from '@/mixins/placeholder'
import Img from '@/components/media/Img.vue'
import Pkg from '@/components/monl/Pkg.vue'
import TagInput from '@/components/tag/TagInput.vue'

export default {
  mixins: [formMixin, modelMixin({ prop: 'pinl' }), placeholderMixin],
  components: {
    Img,
    Pkg,
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
}

.title {
  @apply font-bold;
  @apply mt-1;
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
</style>
