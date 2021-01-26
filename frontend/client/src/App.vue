<template>
  <div id="app">
    <router-view/>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { GameApi } from '@/services/GameApi';
import axios from 'axios';
import { GameManagerApi } from '@/services/GameManagerApi';

const wsHost = (window.location.protocol === 'https:' ? 'wss' : 'ws') + '://' + window.location.host;
const gameApi = new GameApi(axios.create({ baseURL: '/api' }), wsHost);
// todo: move
const gameManagerApi = new GameManagerApi(axios.create({ baseURL: '/api/game-manager' }), wsHost);

@Component({
  provide: {
    gameApi,
    gameManagerApi,
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
