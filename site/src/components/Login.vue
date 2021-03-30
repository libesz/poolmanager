<template>
  <v-main>
    <v-text-field label="Password" v-model="password" type="password"></v-text-field>
    <v-btn @click="tryLogin">Login</v-btn>
  </v-main>
</template>

<script>
  export default {
    name: 'Login',

    data: () => ({
      password: ''
    }),
    methods: {
      tryLogin() {
        fetch('/login', {method: "POST", body: JSON.stringify({password: this.password})})
        .then((result) => {
          console.log(result)
          result.json().then( (decoded) => {
            this.$emit('successfulLogin', decoded.token)
          })
        })
        .catch((err) => console.log(err))
      },
    }
  }
</script>
