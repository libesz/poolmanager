<template>
  <v-container>
    <v-row>
      <v-col class="flex-grow-0">
        <Status :token="token" @loginFailure="$emit('loginFailure')" />
      </v-col>
      <v-col class="flex-grow-0"
        v-for="(controllerConfigSchema, controllerName) in configSchema"
        :key="controllerName"
      >
      <Config :token="token" @loginFailure="$emit('loginFailure')"
        :controllerConfigSchema="controllerConfigSchema"
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
        configSchema: {
          'Temperature controller': {
            Toggles: [
              'Enabled'
            ],
            Ranges: {
              'Temperature': {
                Min: '20',
                Max: '29'
              },
              'Start time': {
                Min: '1',
                Max: '23'
              },
              'End time': {
                Min: '1',
                Max: '23'
              }
            }
          },
          'Pump controller': {
            Toggles: [
              'Enabled'
            ],
            Ranges: {
              'Daily runtime': {
                Min: '1',
                Max: '23'
              }
            }
          }
        }
      }
    },
    props: [
      'token',
    ]
  }
</script>
