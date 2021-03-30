<template>
  <v-main>
    <v-card class="mx-auto" max-width="344" outlined>
      <v-img src="pool_resized.jpg" height="200px"></v-img>
      <v-card-title>
        LOGIN
      </v-card-title>
      <v-list class="transparent">
        <v-list-item>
          <v-text-field label="Password" v-model="password" type="password"></v-text-field>
        </v-list-item>
      </v-list>
      <v-divider></v-divider>
      <v-card-actions>
        <v-btn outlined text @click="tryLogin">
          Login
        </v-btn>
      </v-card-actions>
    </v-card>
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
            if(result.status >= 200 && result.status <= 299){
              this.$emit('successfulLogin', decoded.token)
            } else {
              this.$emit('unsuccessfulLogin', decoded.error)
            } 
          })
        })
        .catch((err) => console.log(err))
      },
    }
  }
</script>
