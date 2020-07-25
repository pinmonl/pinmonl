<template>
 <v-app id="inspire">
    <v-main>
      <v-container
        class="fill-height justify-center"
        fluid
      >
        <v-card class="elevation-12" :style="cardStyles">
          <v-toolbar
            color="primary"
            dark
            flat
          >
            <v-toolbar-title>
              Pinmonl
            </v-toolbar-title>
            <v-spacer></v-spacer>
          </v-toolbar>
          <v-card-text>
            <v-alert type="error" v-if="error">{{ error }}</v-alert>
            <v-form v-model="valid" ref="form">
              <v-text-field
                label="Login"
                name="login"
                v-model="model.login"
                :prepend-icon="mdiAccount"
                type="text"
                :rules="rules.login"
              ></v-text-field>

              <v-text-field
                id="password"
                label="Password"
                name="password"
                v-model="model.password"
                :prepend-icon="mdiLock"
                type="password"
                :rules="rules.password"
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="login" :disabled="!valid" :loading="loading">Login</v-btn>
          </v-card-actions>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { mdiAccount, mdiLock } from '@mdi/js'
import { mapState } from 'vuex'
import { authed } from '@/utils/auth'

export default {
  name: 'login-view',
  data: () => ({
    valid: true,
    model: {
      login: '',
      password: '',
    },
    loading: false,
    rules: {
      login: [
        v => !!v || 'Login is required',
      ],
      password: [
        v => !!v || 'Password is required',
      ],
    },
    error: null,
    mdiAccount,
    mdiLock,
  }),
  computed: {
    ...mapState(['client']),
    cardStyles () {
      const { smAndUp, mdAndUp } = this.$vuetify.breakpoint
      if (mdAndUp) {
        return { width: '480px' }
      }
      if (smAndUp) {
        return { width: '400px' }
      }
      return { 'width': '100%' }
    },
  },
  methods: {
    async login () {
      try {
        this.loading = true
        this.error = null

        const { login, password } = this.model
        const res = await this.client.login(login, password)
        authed(res)
      } catch (e) {
        this.error = e.error || 'Please check your login and password.'
      } finally {
        this.loading = false
      }
    },
  },
}
</script>
