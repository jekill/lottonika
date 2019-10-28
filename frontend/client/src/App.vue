<template>
  <div id="app">
    <router-view/>
  </div>
</template>
<script lang="ts">
  import {Vue, Component} from 'vue-property-decorator';
  import {GameApi} from '@/services/GameApi';
  import axios from 'axios';

  const gameApi = new GameApi(axios.create({baseURL:'/api'}));

  @Component({
    provide: {
      gameApi,
    },
  })
  export default class App extends Vue {
  }
</script>

<style>
  @tailwind base;
  @tailwind components;
  @tailwind utilities;

  html, body, #app {
    height: 100%;
  }
</style>
