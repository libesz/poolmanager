<template>
  <v-main>
    <v-card class="mx-auto" max-width="344" outlined>
      <v-list-item three-line>
        <v-list-item-content>
          <div class="overline mb-1">
            STATUS
          </div>
        </v-list-item-content>
        <v-list-item-content>
          {{status}}
        </v-list-item-content>
      </v-list-item>
      <v-card-actions>
        <v-btn @click="getStatus">Get status</v-btn>
      </v-card-actions>
    </v-card>
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
            if(result.status >= 200 && result.status <= 299){
              result.json()
              .then((decoded) => this.status = decoded)
              .catch((err) => console.log(err))
            } else {
              this.$emit('loginFailure')
            } 
        }).catch((err) => console.log(err))
        .catch((err) => console.log(err))
      }
    }
  }
</script>
