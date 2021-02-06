<template>
  <div class="dashboard-view page-view ">
    <h1>
      Dashboard
      (
      Connection: {{ isWsOnline ? 'Online' : 'Offline' }}
      <template v-if="wsReconnectionTimeout">
        <fa icon="spinner" spin/>
      </template>
      )
    </h1>
    <div class="flex mb-8">
      <q-button @click="handleStartRound">
        Start round
      </q-button>
      <div class="ml-4 pl-4 border-l-2 border-gray-300 flex">
        <q-input-number
            class="w-16"
            v-model="generateCount" type="number" :controls="false"
        />
        <q-button theme="secondary" class="ml-2" @click="handleGenerate">Generate</q-button>
      </div>
      <div class="ml-4 pl-4 border-l-2 border-gray-300">
        <q-button theme="secondary" @click="handleReset">Reset</q-button>
      </div>
    </div>
    <!--    <div class="e-card inline-block">qwe</div>-->

    <div class="flex">
      <div class="mr-8">
        <div
            class="cards-list"
            :style="{width:containerWidth+'px'}"
        >
          <div
              class="cards-list__item"
              :class="{
              'cards-list__item--fail':isFail(card),
              'cards-list__item--win':isWin(card)
          }"
              :style="{
            'top':cardTop(index),
            'left':cardLeft(index)
          }"
              v-for="(card,index) of cards"
              :data-card-number="'card_'+card.number"
              :key="'card_'+card.number"
          >
            <span>{{ card.number }}</span>
          </div>
        </div>
      </div>
      <div cols="2">

        <div
            v-if="[RoundStateEnum.STARTED,RoundStateEnum.FINISHED].includes(roundState)"
            class="round-info"
        >

          <i18n :path="`dashboard.round-is.${roundState}`">
            <b place="number" class="round-info__number">{{ currentRound }}</b>
          </i18n>
          <div
              v-if="roundState===RoundStateEnum.STARTED"
              class="e-card w-32 mt-2 mb-6"
          >
            <transition
                name="fade"
                mode="out-in"
            >
              <span :key="'counter_'+counter">{{ counter }}</span>
            </transition>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script lang="ts">
import { Component, Inject, Vue } from 'vue-property-decorator';
import { CardDto } from '@/models/CardDto';
import { GameManagerApi } from '@/services/GameManagerApi';
import Card from '@/views/Card.vue';
import { RoundState } from '@/types/RoundState';
import { UpdateGameStateMessage } from '@/types/Messages';
import { GameApi } from '@/services/GameApi';

const CARD_WIDTH = 80;
const CARD_HEIGHT = 80;
const GUTTER = 8;
const MAX_IN_ROW = 4;

@Component({})
export default class Dashboard extends Vue {
  @Inject() public gameManagerApi!: GameManagerApi;
  @Inject() public gameApi!: GameApi;
  private cards: CardDto[] = [];
  private isLoading: boolean = false;
  private ws: WebSocket | null = null;
  private isWsOnline: boolean = false;
  private isWsReconnectionEnabled: boolean = true;
  private wsReconnectionTimeout: number | null = null;
  private RoundStateEnum = RoundState;
  private roundState: RoundState = RoundState.NOT_STARTED;
  private currentRound: number = 0;
  private counter: string = '';

  private generateCount: number = 32;

  public created() {
    this.connectWs();
  }

  public destroyed() {
    this.closeWs();
  }

  private handleOnWsClose() {
    this.isWsOnline = false;
    this.tryToReconnect();
  }

  private tryToReconnect() {
    console.log('Reconnection')
    if (this.wsReconnectionTimeout === null && this.isWsReconnectionEnabled) {
      this.wsReconnectionTimeout = window.setTimeout(() => {
        this.wsReconnectionTimeout = null;
        this.connectWs();
      }, 3000);
    }
  }

  private async handleGenerate() {
    for (let i = 0; i < this.generateCount; i++) {
      await this.gameApi.createCard();
    }
  }

  private handleReset() {

  }

  private handleOnWsMessage(ev: MessageEvent) {
    this.isWsOnline = true;
    try {
      if (!ev.data) {
        return;
      }
      const message: UpdateGameStateMessage = JSON.parse(ev.data);
      if (!message?.payload) {
        return;
      }
      this.cards = message.payload.cards || [];
      if (message.payload.round_state !== undefined) {
        this.roundState = message.payload.round_state;
      }
      if (message.payload.counter !== undefined) {
        this.counter = message.payload.counter;
      }
      if (message.payload.current_round) {
        this.currentRound = message.payload.current_round;
      }
      // this.roundState = message.payload?.is_round_started || 0;
    } catch (err) {
      console.error(`FETCHING DATA ERR`, err);
    }
  }

  private cardTop(index: number): string {
    if (this.cards[index].is_closed) {
      return '100%';
    }
    const rowNum = Math.floor(index / MAX_IN_ROW);
    return (rowNum * CARD_HEIGHT + (rowNum) * GUTTER) + 'px';
  }

  private cardLeft(index: number): string {
    const colNum = index % MAX_IN_ROW;
    return (colNum * CARD_WIDTH + (colNum) * GUTTER) + 'px';
  }

  private get containerWidth() {
    return MAX_IN_ROW * CARD_WIDTH + (MAX_IN_ROW - 1) * GUTTER;
  }

  private get containerHeight() {
    const rowsCount = Math.ceil(this.cards.length / MAX_IN_ROW) || 1;
    return CARD_HEIGHT * rowsCount + GUTTER * (rowsCount - 1);
  }

  private connectWs() {
    try {
      this.ws = this.gameManagerApi.wsStateConnect();
      this.ws?.addEventListener('message', this.handleOnWsMessage);
      this.ws?.addEventListener('close', this.handleOnWsClose);
    } catch (e) {
      console.error(e);
      this.tryToReconnect();
    }
  }

  private closeWs() {
    this.isWsReconnectionEnabled = false;
    this.ws?.close();
    this.ws = null;
  }

  private handleStartRound() {
    this.gameManagerApi.startRound();
  }

  private isWin(card: CardDto): boolean {
    return card.is_win;
  }

  private isFail(card: CardDto): boolean {
    return Boolean(card.is_closed && !card.is_win);
  }
}
</script>

<style>
.page-view {
  @apply container mx-auto py-6 px-4;
}

.cards-list {
  /*list-style: none;*/
  /*@apply flex;*/
  /*flex-wrap: wrap;*/
  position: relative;
  padding: 0;
  margin: 0;
  height: 100ch;
  overflow: hidden;
  /*border: 1px solid red;*/
}

.cards-list__item {
  position: absolute;
  @apply block p-6 shadow-xl;
  @apply border-4 rounded border-blue-300;
  @apply text-2xl;
  @apply bg-gray-50;
  width: 80px;
  height: 80px;
  text-align: center;
  transition: top 1s cubic-bezier(0, -0.18, 0.26, -0.43) 3s, background-color 1s ease-in-out;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cards-list__item--transition_top {
  transition: top 0.5s ease-in;
}

.cards-list__item--fail {
  @apply bg-red-500 border-red-600 text-white;
}

.cards-list__item--win {
  @apply bg-green-300 border-green-400;
}

.round-info {
  @apply text-xl;
}

.round-info__number {
  @apply text-2xl;
  color: var(--color-primary-blue-aqua);
}

.e-card {
  min-height: 100px;
  @apply shadow-xl border rounded border-gray-300;
  /*@apply bg-gray-50;*/
  /*@apply bg-blue-800;*/
  background: var(--gradient-secondary);
  @apply text-white text-3xl;
  @apply p-4;
  @apply flex;
  align-items: center;
  justify-content: center;

}

.fade-enter-active, .fade-leave-active {
  transition: opacity .3s ease, color 0.3s ease-in-out;
}

.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
