<template>
  <div class="home">
    <button class="button--green" @click="handleEnterClick">{{ $t('button.enter-game.text') }}</button>
  </div>
</template>

<script lang="ts">
  import {Component, Vue, Inject} from 'vue-property-decorator';
  import {GameApi} from '@/services/GameApi';

  @Component({})
  export default class Home extends Vue {
    @Inject() public gameApi!: GameApi;

    public async handleEnterClick() {
      const cardDto = await this.gameApi.createCard();
      this.$router.push({name: 'card', params: {id: cardDto.id}});
    }
  }
</script>
<style scoped>
  .home {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .button--green {
    @apply bg-green-500 text-white py-2 px-4 rounded;
  }

  .button--green:hover {
    @apply bg-green-400 text-white;
  }

</style>
