<template>
  <div
    class="card-page"
    :class="{
      'card-page--fail': isWin===false,
      'card-page--win': isWin===true
    }"
  >
    <div v-if="!card && !loading">Error</div>
    <div></div>
    <div
      v-if="card"
      class="card"
      :class="{
        'card--fail': isWin===false,
        'card--win': isWin===true
      }"
    >
      <p>{{ $t('card.you-are-in-the-game') }}</p>
      <p class="text-gray-500">{{ $t('card.wait-text') }}</p>
      <div class="card__your-number-text">{{ $t('card.your-number-text') }}</div>
      <div class="card__number">{{ card.number }}</div>
      <div class="counter card__counter">
        {{ counterNumber }}
      </div>
      <div class="card__footer">
        <button
          class="button--red card__stop-button"
          @click="handleStopGame"
        >{{ $t('button.leave-game.text') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  import {Component, Inject, Vue} from 'vue-property-decorator';
  import {GameApi} from '@/services/GameApi';
  import {CardDto} from '@/models/CardDto';

  @Component({})
  export default class Card extends Vue {
    public isWin: boolean | null = null;
    public counterNumber: number = 0;
    public card: CardDto | null = null;
    public loading = true;
    @Inject() public gameApi!: GameApi;

    public async created() {
      const cardId = this.$route.params.id;

      this.loading = true;
      try {
        this.card = await this.gameApi.getCard(cardId) || null;
      }catch (e) {
      }finally {
        this.loading = false;
      }

      // this.card = {
      //   id: '333',
      //   number: 333,
      // };
      // setInterval(() => {
      //   this.counterNumber++;
      // }, 400);
    }

    public async handleStopGame() {
      if (confirm(String(this.$t('stop-game.are-you-sure.confirmation-text')))) {
        if (this.card) {
          await this.gameApi.stopGame(this.card.id);
        }
        await this.$router.push({name: 'home'});
      }
    }
  }
</script>

<style scoped>
  .card-page {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: center;

    height: 100%;
  }

  .card-page--fail {
    @apply bg-red-600;
  }

  .card-page--win {
    @apply bg-green-400;
  }

  .card {
    @apply relative;
    @apply max-w-sm shadow-lg overflow-hidden;
    @apply rounded border-2;
    @apply px-6 py-4 m-6;
    width: 100%;
    text-align: center;
    height: 80%;
  }

  .card__number {
    @apply font-extrabold text-blue-500;
    font-size: 120px;
  }

  .card--fail {
    @apply bg-red-500;
  }

  .card--win {
    @apply bg-green-300;
  }

  .card__footer {
    text-align: center;
    @apply absolute bottom-0 left-0 mb-4 w-full;
  }

  .counter {
    @apply text-6xl;
  }

  .button--red {
    @apply bg-red-500 text-white py-2 px-4 rounded;
  }

  .button--red:hover {
    @apply bg-red-400 text-white ;
  }
</style>
