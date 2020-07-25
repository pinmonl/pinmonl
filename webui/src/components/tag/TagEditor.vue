<template>
  <v-dialog
    :value="show"
    @input="$emit('update:show', $event)"
    max-width="600"
    scrollable
    :fullscreen="$vuetify.breakpoint.mobile"
  >
    <v-card>
      <v-card-title>
        {{ isNew ? 'Create' : 'Edit' }}
      </v-card-title>
      <v-card-text>
        <v-form v-if="model" v-model="valid" ref="form">
          <v-text-field
            :label="isNew ? 'Name' : 'Rename'"
            autofocus
            prefix="/"
            v-model="modelName"
            v-bind="modelNameProps"
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn text :disabled="loading" @click="$emit('update:show', false)">
          Cancel
        </v-btn>
        <v-btn text color="error" @click="remove">
          Delete
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn text color="primary" @click="save">
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex'
import cloneDeep from 'lodash.clonedeep'

export default {
  name: 'tag-editor',
  components: {
  },
  props: {
    show: { type: Boolean },
    value: { type: Object },
    autofocus: { type: Boolean },
  },
  data () {
    return {
      valid: true,
      loading: false,
      updateStatus: '',
      updateMessage: '',
      model: cloneDeep(this.value),
      isModelNameUsed: false,
    }
  },
  computed: {
    ...mapState(['client']),
    isNew () {
      return !(this.value || {}).id
    },
    tagName () {
      return (this.value || {}).name
    },
    modelName: {
      get () {
        return (this.model || {}).name
      },
      set (val) {
        if (!this.model) {
          return
        }
        this.model.name = val
      },
    },
    modelNameHint () {
      if (this.isNew) {
        return
      }
      if (this.modelName == this.tagName) {
        return
      }
      return `Update tags with prefix "/${this.tagName}" to "/${this.modelName}".`
    },
    nameRules () {
      return [
        v => v.length >= 1 || 'At least has one character',
        v => !/^\//.test(v) || 'Incorrect pattern',
        v => !/\/$/.test(v) || 'Incorrect pattern',
      ]
    },
    modelNameProps () {
      const props = {
        rules: this.nameRules,
        hint: this.modelNameHint,
        persistentHint: true,
      }

      if (this.isModelNameUsed) {
        props['error'] = true
        props['error-messages'] = this.modelNameErrorMessage
      }

      return props
    },
    modelNameErrorMessage () {
      if (this.isModelNameUsed) {
        return 'Name is used'
      }
      return null
    },
  },
  methods: {
    async search (tagName) {
      const res = await this.client.listTags({ q: tagName })
      return res.data
    },
    async save () {
      this.$refs.form.validate()
      if (!this.valid) {
        return
      }

      try {
        this.loading = true

        const available = await this.checkNameAvailability()
        if (!available) {
          this.isModelNameUsed = true
          return
        }

        const { id } = this.model
        const res = this.isNew
          ? await this.client.createTag(this.model)
          : await this.client.updateTag(id, this.model)

        this.model = res
        this.$emit('change', res)
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    async remove () {
      try {
        this.loading = true

        await this.client.deleteTag(this.value.id)

        this.$emit('delete', this.value)
      } catch (e) {
        console.log(e)
        //
      } finally {
        this.loading = false
      }
    },
    async checkNameAvailability () {
      if (this.tagName == this.modelName) {
        return true
      }

      const found = await this.search(`/${this.modelName}/`)
      return found.length == 0
    },
  },
  watch: {
    value (val) {
      this.model = cloneDeep(val)
    },
    show (val) {
      if (!val) {
        return
      }

      this.$nextTick(() => {
        this.$refs.form.resetValidation()
      })
    },
    modelName () {
      this.isModelNameUsed = false
    },
  },
}
</script>
