<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
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
      <Main v-if="token" :token="token" />
      <Login v-else @successfulLogin="successfulLogin" />
      <p>
        <span>Logged </span>
        <span v-if="token">in, token: {{token}}</span>
        <span v-else>out</span>
      </p>
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

  data: () => {
    return {
      token: '',
      status: ''
    }
  },
  mounted() {
    if (localStorage.token) {
      this.token = localStorage.token;
    }
  },
  /*watch: {
    token(token) {
      localStorage.token = token;
    }
  },*/
  methods: {
    logout() {
      localStorage.token = null
      this.token = null
    },
    successfulLogin(token) {
      localStorage.token = token
      this.token = token
    }
  }
};
</script>
