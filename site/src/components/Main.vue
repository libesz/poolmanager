<template>
  <v-main>
    <Status :token="token" @loginFailure="$emit('loginFailure')" />
  </v-main>
</template>

<script>
  import Status from './Status'

  export default {
    name: 'Main',
    components: {
      Status
    },
    props: [
      'token',
    ]
  }
</script>
