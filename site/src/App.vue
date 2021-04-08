<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <v-img
          alt="Pool Manager Logo"
          class="shrink mr-2"
          contain
          src="/favicon.png"
          transition="scale-transition"
          width="40"
        />

        <span class="display-1">Pool manager</span>
      </div>

      <v-spacer></v-spacer>

      <v-btn v-if="token" @click="logout" text>
        <span class="mr-2">Logout</span>
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <Main v-if="token" :token="token" @loginFailure="loginFailure" />
      <Login v-else @successfulLogin="successfulLogin" @loginFailure="loginFailure" />
      <div class="text-center">
        <v-snackbar v-model="snackbar" :timeout="snackbarTimeout">
          {{ snackbarText }}
          <template v-slot:action="{ attrs }">
            <v-btn color="blue" text v-bind="attrs" @click="snackbar = false" >
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </template>
        </v-snackbar>
      </div>
    </v-main>
  </v-app>
</template>

<script>

import Main from './components/Main'
import Login from './components/Login'

export default {
  name: 'App',

  components: {
    Main,
    Login
  },

  data: () => ({
      token: '',

      snackbar: false,
      snackbarText: '',
      snackbarTimeout: 4000,
    }
  ),
  created() {
    console.log(localStorage.token)
    if (localStorage.token && localStorage.token != 'null') {
      this.token = localStorage.token;
    }
  },
  watch: {
    token(token) {
      localStorage.token = token;
    }
  },
  methods: {
    logout() {
      this.token = ''
    },
    successfulLogin(token) {
      this.token = token
    },
    loginFailure(error) {
      if(!error) {
        error = "Unexpected error. Please log in again."
      }
      this.snackbarText = error
      this.snackbar = true
      this.logout()
    }
  }
};
</script>
