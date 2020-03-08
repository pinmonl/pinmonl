<template>
  <div :class="$style.container">
    <template v-if="loading">
      <Placeholder>
        <div :class="[$style.phRow, $style.phLine, $style.phTitle]" />
        <div :class="[$style.phRow, $style.phLine]" />
        <div :class="[$style.phRow, $style.phLine]" />
        <div :class="[$style.phRow, $style.phLine, $style.phLineShort]" />
      </Placeholder>
    </template>

    <template v-else-if="editable">
      <InputGroup>
        <Label>Name</Label>
        <Input v-model="model.name" />
      </InputGroup>
      <InputGroup>
        <Label>Must tags</Label>
        <TagInput v-model="mustTags" />
      </InputGroup>
      <InputGroup>
        <Label>Any tags</Label>
        <TagInput v-model="anyTags" />
      </InputGroup>
      <InputGroup>
        <Label>Description</Label>
        <Textarea v-model="model.description" />
      </InputGroup>
      <InputGroup>
        <Label>Readme</Label>
        <Textarea v-model="model.readme" />
      </InputGroup>
    </template>

    <template v-else>
      <div :class="$style.name" v-text="model.name" />
      <TagInput :class="$style.mustTags" :value="mustTags" disabled noStyle />
      <TagInput :class="$style.anyTags" :value="anyTags" disabled noStyle />
      <div :class="$style.description" v-text="model.description" />
      <div :class="$style.readme" v-text="model.readme" />
    </template>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import modelMixin from '@/mixins/model'
import placeholderMixin from '@/mixins/placeholder'
import TagInput from '@/components/tag/TagInput.vue'

export default {
  mixins: [formMixin, modelMixin({ prop: 'share' }), placeholderMixin],
  components: {
    TagInput,
  },
  props: {
    share: {
      type: Object,
      required: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    model () {
      const gt = this.getter
      const st = this.setter
      return {
        $data: this.share,
        get name () { return gt('name') },
        set name (v) { st('name', v) },
        get description () { return gt('description') },
        set description (v) { st('description', v) },
        get readme () { return gt('readme') },
        set readme (v) { st('readme', v) },
        get mustTags () { return gt('mustTags') },
        set mustTags (v) { st('mustTags', v) },
        get anyTags () { return gt('anyTags') },
        set anyTags (v) { st('anyTags', v) },
      }
    },
    mustTags: {
      get () {
        const tags = this.$store.getters['tag/tags']
        return this.$store.getters['tag/getByName'](tags, this.model.mustTags)
      },
      set (tags) {
        const tagNames = this.$store.getters['tag/mapName'](tags)
        this.model.mustTags = tagNames
      },
    },
    anyTags: {
      get () {
        const tags = this.$store.getters['tag/tags']
        return this.$store.getters['tag/getByName'](tags, this.model.anyTags)
      },
      set (tags) {
        const tagNames = this.$store.getters['tag/mapName'](tags)
        this.model.anyTags = tagNames
      },
    },
  },
}
</script>

<style lang="scss" module>
.container {
  @apply p-4;
}

.name,
.description,
.mustTags,
.anyTags,
.readme {
  @apply mb-2;
}

.phRow {
  @apply mb-2;
}

.phLine {
  height: 1.2rem;
}

.phTitle {
  width: 50%;
}

.phLineShort {
  width: 30%;
}
</style>
