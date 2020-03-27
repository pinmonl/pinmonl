<template>
  <div :class="$style.container">
    <Header noAccount>
      <h2 :class="$style.title">My account</h2>
    </Header>
    <div :class="$style.panel">
      <Form>
        <InputGroup>
          <Label>Login</Label>
          <Input v-model="model.login" />
        </InputGroup>
        <InputGroup>
          <Label>Name</Label>
          <Input v-model="model.name" />
        </InputGroup>
        <InputGroup>
          <Label>Change password</Label>
          <Input type="password" v-model="password" />
        </InputGroup>
        <div :class="$style.buttonContainer">
          <Button :class="$style.button" @click="handleSave">
            Update
          </Button>
          <Button :class="$style.button" @click="handleLogout" light>
            Log out
          </Button>
        </div>
      </Form>
    </div>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import { validationMixin } from 'vuelidate'
import { required, alphaNum, minLength } from 'vuelidate/lib/validators'
import Header from '@/components/app/Header.vue'

export default {
  mixins: [formMixin, validationMixin],
  components: {
    Header,
  },
  data () {
    return {
      inputs: null,
      password: null,
    }
  },
  created () {
    this.inputs = this.getDefaultValue()
  },
  computed: {
    user () {
      return this.$store.getters.user
    },
    model () {
      const gt = this.getter
      const st = this.setter
      return {
        $data: this.user,
        get name () { return gt('name') },
        set name (v) { return st('name', v) },
        get login () { return gt('login') },
        set login (v) { return st('login', v) },
      }
    },
  },
  methods: {
    getDefaultValue () {
      return { ...this.user }
    },
    getter (key) {
      return this.inputs[key]
    },
    setter (key, value) {
      this.inputs = { ...this.inputs, [key]: value }
    },
    async handleSave () {
      await this.$store.dispatch('updateMe', this.inputs)
      this.inputs = this.getDefaultValue()
      this.password = null
    },
    async handleLogout () {
      await this.$store.dispatch('logout')
      this.$router.push({ name: 'login' })
    },
  },
  validations () {
    const v = {
      login: { required, format: alphaNum },
      name: { required },
      password: { required, minLength: minLength(6) },
    }
    return { inputs: v }
  },
}
</script>

<style lang="scss" module>
.container {
}

.title {
  @apply flex-grow;
}

.panel {
  max-width: 400px;
  @apply mx-auto;
  @apply p-4;
}

.input {
  @apply mt-2;
}

.buttonContainer {
  @apply mt-6;
  @apply flex;
  @apply justify-between;
  @apply flex-wrap;
}

.button {
  @apply w-full;
  @apply mb-4;

  @screen md {
    width: 48%;
  }
}
</style>
