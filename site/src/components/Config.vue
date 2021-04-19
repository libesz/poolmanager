<template>
  <v-card width="344" height="100%" outlined>
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
              v-for="(toggle, index) in controllerConfigSchema.toggles"
              :key="toggle.name+index"
              >
                <td>{{ toggle.name }}</td>
                <td><v-switch @change="switchChange(controllerName, toggle.name, $event)" v-model="controllerConfigData.toggles[toggle.name]"></v-switch></td>
              </tr>
              <tr
              v-for="(range, index) in controllerConfigSchema.ranges"
              :key="range.name + index"
              >
                <td>{{ range.name }}</td>
                <td width="70%"><v-slider :hint="range.name" :min="range.min" :max="range.max" :step="range.step" @change="sliderChange(controllerName, range.name, $event)" thumb-label="always" :value="controllerConfigData.ranges[range.name]"></v-slider></td>
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
      controllerConfigData: Object,
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
      },
      switchChange(controllerName, toggleName, value) {
        console.log(`Toggle change: ${controllerName} ${toggleName} ${value}`)
        fetch('/api/config', {
          method: 'POST',
          headers: {'Authorization': 'Bearer ' + this.$props.token},
          body: JSON.stringify({'controller':controllerName, 'type':'toggle', 'key':toggleName, 'value':value.toString()})
        })
        .then((result) => {
          result.json()
          .then((decoded) => {
            if(result.status >= 200 && result.status <= 299){
              console.log('ok')
            } else {
              this.$emit('configError', decoded)
            }
          })
          .catch(() => this.$emit('loginFailure'))
        }).catch((err) => console.log(err))
        .catch((err) => console.log(err))
      },
      sliderChange(controllerName, rangeName, value) {
        console.log(`Slider change: ${controllerName} ${rangeName} ${value}`)
        fetch('/api/config', {
          method: 'POST',
          headers: {'Authorization': 'Bearer ' + this.$props.token},
          body: JSON.stringify({'controller':controllerName, 'type':'range', 'key':rangeName, 'value':value.toString()})
        })
        .then((result) => {
          result.json()
          .then((decoded) => {
            if(result.status >= 200 && result.status <= 299){
              console.log('ok')
            } else {
              this.$emit('configError', decoded.error)
            }
          })
          .catch(() => this.$emit('loginFailure'))
        }).catch((err) => console.log(err))
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
