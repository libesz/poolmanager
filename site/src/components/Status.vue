<template>
  <v-card width="344" height="100%" outlined>
    <v-card-title>
      STATUS
      <v-spacer></v-spacer>
      <v-icon>mdi-coolant-temperature</v-icon>
    </v-card-title>
    <v-list class="transparent">
      <v-list-item>
        <v-simple-table>
          <template v-slot:default>
            <tbody>
              <tr
              v-for="(inputValue, name) in status.inputs"
              :key="inputValue"
              >
                <td>{{ name }}</td>
                <td>{{ inputValue }}</td>
              </tr>
              <tr
              v-for="(outputValue, name) in status.outputs"
              :key="outputValue"
              >
                <td>{{ name }}</td>
                <td><v-icon v-if="outputValue" color="blue">mdi-brightness-1</v-icon><v-icon v-else>mdi-brightness-1</v-icon></td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-list-item>
    </v-list>
    <!--v-card-actions>
      <v-btn @click="getStatus">Update</v-btn>
    </v-card-actions-->
  </v-card>
</template>

<script>
  export default {
    name: 'Status',

    props: [
      'token',
    ],
    data: () => {
      return {
        status: ''
      }
    },
    created() {
      this.getStatus()
      setInterval(this.getStatus, 5000)
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
<style lang="scss">  
  tbody {
     tr:hover {
        background-color: transparent !important;
     }
  }
</style>
