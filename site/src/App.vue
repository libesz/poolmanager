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
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
          transition="scale-transition"
          width="40"
        />

        <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        />
      </div>

      <v-spacer></v-spacer>

      <v-btn
        href="https://github.com/vuetifyjs/vuetify/releases/latest"
        target="_blank"
        text
      >
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <div v-if="token">
      <Main :token="token" />
      <v-btn @click="logout">Logout</v-btn>
      </div>
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
