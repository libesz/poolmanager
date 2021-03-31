<template>
  <v-card class="mx-auto" width="344" outlined>
    <v-card-title>
      {{ controllerName }}
      <v-spacer></v-spacer>
      <v-icon>mdi-cog-outline</v-icon>
    </v-card-title>
    <v-list class="transparent">
      <v-list-item>
        <v-layout child-flex>
        <v-simple-table>
            <tbody>
              <tr
              v-for="(toggle, index) in controllerConfigSchema.Toggles"
              :key="toggle+index"
              >
                <td>{{ toggle }}</td>
                <td>{{ toggle }}</td>
              </tr>
              <tr
              v-for="(rangeProperties, rangeName, index) in controllerConfigSchema.Ranges"
              :key="rangeName + index"
              >
                <td>{{ rangeName }}</td>
                <td>{{ rangeProperties.Min}} {{rangeProperties.Max}}</td>
              </tr>
            </tbody>
          </v-simple-table>
        </v-layout>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script>
  export default {
    name: 'Config',

    props: {
      token: String,
      controllerConfigSchema: Object,
      controllerName: String
    },
    data: () => {
      return {


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
<style lang="scss">  
  tbody {
     tr:hover {
        background-color: transparent !important;
     }
  }
</style>
