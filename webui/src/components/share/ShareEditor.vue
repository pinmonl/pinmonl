<template>
  <v-dialog v-model="vShow" scrollable maxWidth="600">
    <v-stepper v-model="step">
      <v-stepper-header>
        <v-stepper-step step="1" editable>Share</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step step="2" editable>Select tags</v-stepper-step>
      </v-stepper-header>

      <v-stepper-items>
        <v-stepper-content step="1">
          <v-card>
            <v-card-text>
              <v-form ref="form1" v-model="valid1">
                <v-text-field
                  v-model="internalValue.slug"
                  label="Slug"
                  :rules="slugRules"
                  :autofocus="autofocus"
                ></v-text-field>
                <v-text-field
                  v-model="internalValue.name"
                  label="Name"
                ></v-text-field>
                <tag-input
                  v-model="internalValue.tags"
                  label="Must tags"
                  hint="Pins having all of these tags will be included in share."
                  persistent-hint
                  :rules="tagsRules"
                ></tag-input>
              </v-form>
            </v-card-text>

            <v-card-actions>
              <v-btn text @click="vShow = false">Cancel</v-btn>
              <v-btn v-if="!isNew" text @click="() => {}" color="error">Delete</v-btn>
              <v-spacer></v-spacer>
              <v-btn @click="() => {}" color="primary">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-stepper-content>

        <v-stepper-content step="2">
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-dialog>
</template>

<script>
import TagInput from '@/components/tag/TagInput'

export default {
  name: 'share-editor',
  components: {
    TagInput,
  },
  props: {
    value: { type: Object },
    show: { type: Boolean },
    autofocus: { type: Boolean },
  },
  data: () => ({
    loading: false,
    valid1: false,
    // internalValue: null,
    internalValue: { name: '', slug: '', tags: [] },
    step: 1,

    slugRules: [
      v => !!v || 'Required',
      v => /^[a-z0-9-]+$/i.test(v) || 'Only a-z, 0-9 and - are allowed',
    ],
    tagsRules: [
      v => !!v && v.length > 0 || 'Required',
    ],
  }),
  computed: {
    vShow: {
      get () { return this.show },
      set (v) { this.$emit('update:show', v) },
    },
    isNew () { return !this.internalValue.id },
  },
  methods: {
  },
  watch: {
    value (val) {
      this.internalValue = val
    },
  },
}
</script>
