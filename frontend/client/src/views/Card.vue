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
      <div v-if="isWsOpen">Connected</div>
      <div v-else="isWsOpen">Offline</div>

      <p v-if="roundState === RoundStateEnum.NOT_STARTED">{{ $t('card.you-are-in-the-game') }}</p>

      <p
          class="text-gray-500"
          v-if="roundState===RoundStateEnum.STARTED"
      >
        {{ $t('card.round-started') }}
      </p>

      <p class="text-gray-500" v-if="roundState===RoundStateEnum.NOT_STARTED">{{ $t('card.wait-text') }}</p>
      <p v-else class="round-info">
        <i18n :path="`dashboard.round-is.${roundState}`">
          <b place="number" class="round-info__number">{{ currentRound }}</b>
        </i18n>
      </p>

      <div class="card__your-number-text">{{ $t('card.your-number-text') }}</div>
      <div
          class="card__number"
          :class="{'card__number--small':isCounterShown}"
      >{{ card.number }}
      </div>

      <transition
          v-if="isCounterShown"
          name="number-change"
          mode="out-in"
      >
        <div class="counter card__counter" v-bind:key="'counter-number_'+counterNumber">
          {{ counterNumber }}
        </div>
      </transition>
      <div class="card__footer">
        <!--
        <button
            class="button--red card__stop-button"
            @click="handleStopGame"
        >
          {{ $t('button.leave-game.text') }}
        </button>
        -->
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Inject, Vue } from 'vue-property-decorator';
import { GameApi } from '@/services/GameApi';
import { CardDto } from '@/models/CardDto';
import { CommunicationMessages, RefreshMessage } from '@/types/Messages';
import { RoundState } from '@/types/RoundState';

@Component({})
export default class Card extends Vue {
  public counterNumber: string = '';
  public card: CardDto | null = null;
  public roundState: RoundState = RoundState.NOT_STARTED;
  public RoundStateEnum = RoundState;
  public loading = true;
  public ws: WebSocket | null = null;
  public wsReadyState: number = WebSocket.CLOSED;
  @Inject() public gameApi!: GameApi;

  public destroyed() {
    this.ws?.close();
  }

  public async created() {
    const cardId = this.$route.params.id;
    // const inc = () => setTimeout(() => {
    //   this.counterNumber++;
    //   inc();
    // }, 5000);
    // inc();

    this.loading = true;
    try {
      this.card = await this.gameApi.getCard(cardId) || null;
      // debugger;
      if (this.card?.is_closed === false) {
        this.connectWs();
      }
      // (<any>window).__ws = this.ws;
    } catch (e) {
      console.error(e);
    } finally {
      this.loading = false;
    }
  }

  private connectWs() {
    if (!this.card) {
      return;
    }
    this.ws = this.gameApi.wsConnect(this.card.id);
    this.ws.addEventListener('message', this.onServerMessage);
    this.ws.addEventListener('open', () => {
      this.wsReadyState = this.ws?.readyState;
      this.ws?.send(JSON.stringify({ id: this.card.id, type: 'hello', payload: 'hello' }));
    });
    this.ws.addEventListener('close', () => {
      this.wsReadyState = this.ws?.readyState;
      console.log('WS closed');
    });
  }

  public get isCounterShown(): boolean {
    return this.roundState === RoundState.STARTED;
  }

  public get isWin(): boolean | null {
    if (this.card?.is_closed === true) {
      return Boolean(this.card.is_win);
    }

    return this.card?.is_win === true || null;
  }

  public onServerMessage(message: MessageEvent<string>) {
    this.wsReadyState = this.ws?.readyState;
    let data: CommunicationMessages;
    try {
      data = JSON.parse(message.data) as CommunicationMessages;
    } catch (err) {
      console.error("Parse data error", err);
      return;
    }

    if (!data) {
      return;
    }

    if (data.payload?.current_round !== undefined) {
      this.currentRound = data.payload.current_round;
    }

    if (data.payload?.round_state !== undefined) {
      this.roundState = data.payload.round_state;
    }

    if (data.payload?.card !== undefined) {
      this.card = (data as RefreshMessage).payload.card;
    }

    if (data.payload?.counter !== undefined) {
      this.counterNumber = data.payload.counter;
    }
    console.log(message);
  }

  public get isWsOpen() {
    return this.ws !== null && this.wsReadyState === WebSocket.OPEN;
  }

  public async handleStopGame() {
    if (confirm(String(this.$t('stop-game.are-you-sure.confirmation-text')))) {
      if (this.card) {
        await this.gameApi.stopGame(this.card.id);
      }
      await this.$router.push({ name: 'home' });
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
    transition: background-color 1s ease-in-out;

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
    transition: background-color 1s ease-in-out;
  }

  .card__number {
    @apply font-extrabold text-blue-500;
    @apply text-9xl;
    transition: font-size 0.3s ease-in-out;
  }

  .card__number--small {
    @apply text-6xl;
  }

  .card__counter {
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
    @apply text-9xl;
  }

  .button--red {
    @apply bg-red-500 text-white py-2 px-4 rounded;
  }

  .button--red:hover {
    @apply bg-red-400 text-white ;
  }

  .number-change-enter-active, .number-change-leave-active {
    transition: opacity .3s ease, font-size 0.3s ease-in-out, color 0.3s ease-in-out;
  }

  .number-change-enter, .number-change-leave-to {
    opacity: 0;
  }

  .number-change-enter {
    font-size: 0.4rem;
  }

  .number-change-leave-to {
    /*font-size: 2rem;*/
  }
</style>
