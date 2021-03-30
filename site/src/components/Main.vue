<template>
  <v-main>
  <v-btn @click="getStatus">Get status</v-btn>
  <span>{{status}}</span>
  </v-main>
</template>

<script>
  export default {
    name: 'Main',

    props: [
      'token',
    ],
    data: () => {
      return {
        status: ''
      }
    },
    methods: {
      getStatus() {
        fetch('/api/status', {headers: {'Authorization': 'Bearer ' + this.$props.token}})
        .then((result) => {
          result.json()
          .then((decoded) => this.status = decoded)
          .catch((err) => console.log(err))
        }).catch((err) => console.log(err))
        .catch((err) => console.log(err))
      }
    }
  }
</script>
