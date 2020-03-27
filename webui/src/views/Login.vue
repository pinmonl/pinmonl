<template>
  <div :class="$style.root">
    <div :class="$style.container">
      <Brand tag="h1" />
      <h3>
        <template v-if="isSignup">
          Sign up
        </template>
        <template v-if="!isSignup">
          Login
        </template>
      </h3>
      <Form :class="$style.form">
        <InputGroup v-if="isSignup">
          <Input placeholder="Name" v-model="model.name" />
        </InputGroup>
        <InputGroup>
          <Input placeholder="Username" v-model="model.login" />
          <template #errors v-if="$v.model.login.$error">
            <Errors>
              <p v-if="!$v.model.login.required">Please enter your username.</p>
              <p v-if="!$v.model.login.format">Please enter correct username.</p>
            </Errors>
          </template>
        </InputGroup>
        <InputGroup>
          <Input type="password" placeholder="Password" v-model="model.password" />
          <template #errors v-if="$v.model.password.$error">
            <Errors>
              <p v-if="!$v.model.password.required">Please enter your password.</p>
              <p v-if="!$v.model.password.minLength">Please enter your password.</p>
            </Errors>
          </template>
        </InputGroup>
        <Button :class="$style.button" @click="submit">
          <template v-if="isSignup">
            Sign up
          </template>
          <template v-if="!isSignup">
            Login
          </template>
        </Button>
      </Form>
      <div :class="$style.footer">
        <template v-if="isSignup">
          Already have an account?
          <Anchor :to="{ name: 'login' }" color>Login</Anchor>
        </template>
        <template v-if="!isSignup">
          Do not have an account yet?
          <Anchor :to="{ name: 'signup' }" color>Sign up</Anchor>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import formMixin from '@/mixins/form'
import { mapActions } from 'vuex'
import { validationMixin } from 'vuelidate'
import { required, alphaNum, minLength } from 'vuelidate/lib/validators'
import Brand from '@/components/app/Brand.vue'

export default {
  mixins: [formMixin, validationMixin],
  components: {
    Brand,
  },
  props: {
    new: {
      type: Boolean,
      default: false,
    },
  },
  data () {
    return {
      model: null,
    }
  },
  created () {
    this.model = this.getDefaultValue()
  },
  computed: {
    isSignup () {
      return this.new
    },
  },
  methods: {
    async submit () {
      this.$v.$touch()
      if (this.$v.$invalid) {
        return
      }

      try {
        if (this.isSignup) {
          await this.signup({
            login: this.model.login,
            password: this.model.password,
            name: this.model.name,
          })
        } else {
          await this.login({
            login: this.model.login,
            password: this.model.password,
          })
        }
        this.$router.push({ name: 'home' })
      } catch (e) {
        //
      }
    },
    clear () {
      this.model = this.getDefaultValue()
    },
    getDefaultValue () {
      return {
        login: '',
        name: '',
        password: '',
      }
    },
    ...mapActions([
      'login',
      'signup',
    ]),
  },
  validations () {
    const v = {
      login: { required, format: alphaNum },
      password: { required, minLength: minLength(6) },
    }

    if (this.isSignup) {
      v.login.format = alphaNum
      v.name = { required }
    }
    return { model: v }
  },
}
</script>

<style lang="scss" module>
.root {
  @apply h-full;
  @apply w-full;
}

.container {
  @apply mt-20;
  @apply mx-auto;
  @apply p-8;

  @screen md {
    max-width: 400px;
  }

  h1 {
    @apply mb-6;
    @apply text-3xl;
  }

  h3 {
    @apply mb-1;
  }
}

.button {
  @apply w-full;
}

.footer {
  @apply py-4;
  @apply text-sm;
}
</style>
