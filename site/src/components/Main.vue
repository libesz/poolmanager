<template>
  <v-container v-if="configSchema !== null">
    <v-row>
      <v-col class="flex-grow-0">
        <Status :token="token" @loginFailure="$emit('loginFailure')" />
      </v-col>
      <v-col class="flex-grow-0"
        v-for="(controllerConfigSchema, controllerName) in configSchema"
        :key="controllerName"
      >
        <Config :token="token" @loginFailure="$emit('loginFailure')" @configError="$emit('configError', $event)"
          :controllerConfigSchema="controllerConfigSchema"
          :controllerConfigData="initialConfigData[controllerName]"
          :controllerName="controllerName" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  import Status from './Status'
  import Config from './Config'

  export default {
    name: 'Main',
    components: {
      Status,
      Config
    },
    data: () => {
      return {
        configSchema: null,
        initialConfigData: null
      }
    },
    props: [
      'token',
    ],
    created() {
      fetch('api/config', {headers: {'Authorization': 'Bearer ' + this.$props.token}})
      .then((result) => {
        if(result.status >= 200 && result.status <= 299){
          result.json()
          .then((decoded) => {
            this.configSchema = decoded.schema
            this.initialConfigData = decoded.values
            console.log(decoded)
          })
          .catch((err) => console.log(err))
        } else {
          this.$emit('loginFailure')
        } 
      }).catch((err) => console.log(err))
      .catch((err) => console.log(err))
    }
  }
</script>
