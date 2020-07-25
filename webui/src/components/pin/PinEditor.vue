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
        <v-form ref="form" v-if="model" v-model="valid">
          <image-input
            :image-id="model.imageId"
            width="100"
            height="100"
            :data-uri-src="imageDataUri"
            :disabled="loading"
            @change="(file) => image = file"
          ></image-input>

          <v-text-field
            label="URL"
            v-model="model.url"
            :rules="urlRules"
            :disabled="loading"
            :autofocus="autofocus"
          >
            <template #append>
              <v-btn text :loading="loading" icon v-if="model.url" @click="fillByCard" tabindex="-1">
                <v-icon>{{ mdiReload }}</v-icon>
              </v-btn>
            </template>
          </v-text-field>

          <v-text-field
            label="Title"
            v-model="model.title"
            :rules="titleRules"
            :disabled="loading"
          ></v-text-field>

          <v-text-field
            label="Description"
            v-model="model.description"
            :disabled="loading"
          ></v-text-field>

          <tag-input
            label="Tags"
            v-model="model.tags"
            :disabled="loading"
            multiple
          ></tag-input>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn text :disabled="loading" @click="$emit('update:show', false)">
          Cancel
        </v-btn>
        <v-btn text :loading="loading" color="error" @click="remove">
          Delete
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn text :loading="loading" color="primary" @click="save">
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import ImageInput from '@/components/form/ImageInput'
import TagInput from '@/components/tag/TagInput'
import { mapState } from 'vuex'
import { base64ToBlob } from '@/utils/datauri'
import { mdiReload } from '@mdi/js'

export default {
  name: 'pin-editor',
  components: {
    ImageInput,
    TagInput,
  },
  props: {
    show: { type: Boolean },
    value: { type: Object },
    autofocus: { type: Boolean },
  },
  data () {
    return {
      image: null,
      imageDataUri: '',
      loading: false,
      model: this.value,

      valid: true,

      urlRules: [
        v => !!v || 'Required',
        v => /^https?:\/\//.test(v) || 'Invalid format',
      ],
      titleRules: [
        v => !!v || 'Required',
      ],

      mdiReload,
    }
  },
  computed: {
    ...mapState(['client']),
    isNew () {
      return !(this.value || {}).id
    },
  },
  methods: {
    fieldModel (key) {
      return {
        value: this.model[key],
        event: (v) => this.model[key] = v,
      }
    },
    async save () {
      this.$refs.form.validate()
      if (!this.valid) {
        return
      }

      try {
        this.loading = true

        const res = this.isNew
          ? await this.client.createPin(this.model)
          : await this.client.updatePin(this.model.id, this.model)
        if (this.image) {
          const imgRes = await this.client.uploadPinImage(res.id, this.image)
          res.imageId = imgRes.id
        }

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

        await this.client.deletePin(this.value.id)

        this.$emit('delete', this.value)
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
    async fillByCard () {
      try {
        this.loading = true
        const card = await this.client.card(this.value.url)
        this.value.title = card.title
        this.value.description = card.description

        if (card.imageData) {
          this.imageDataUri = card.imageData
          this.image = base64ToBlob(card.imageData, 'image/png')
        } else {
          this.imageDataUri = null
          this.image = null
        }
      } catch (e) {
        //
      } finally {
        this.loading = false
      }
    },
  },
  watch: {
    value () {
      this.model = this.value
      this.image = null
      this.imageDataUri = null
    },
    show (val) {
      if (!val) {
        return
      }

      this.$nextTick(() => {
        this.$refs.form.resetValidation()
      })
    },
  },
}
</script>
