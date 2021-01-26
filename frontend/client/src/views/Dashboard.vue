<template>
  <div class="dashboard-view page-view">
    <h1>Dashboard</h1>
    <div class="flex">
      <button @click="handleStartRound">Start round</button>
      <div class="ml-4">
        <input v-model="generateCount" type="number" class="border w-12">
        <button class="ml-2" @click="handleGenerate">Generate</button>
      </div>
      <div class="ml-4">
        <button @click="handleReset">Reset</button>
      </div>
    </div>
    <div v-if="roundState>0">
      A round is started
      <div class="e-card">
        {{ counter }}
      </div>
    </div>
    {{ containerWidth }}
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
      >
        {{ card.number }}
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
  private roundState: RoundState = RoundState.IS_NOT_STARTED;
  private counter: number;

  private generateCount: number = 32;

  public created() {
    this.connectWs();
    this.ws?.addEventListener('message', this.handleOnWsMessage);
  }

  public destroyed() {
    this.closeWs();
  }

  private async handleGenerate() {
    for (let i = 0; i < this.generateCount; i++) {
      await this.gameApi.createCard();
    }
  }

  private handleReset() {

  }

  private handleOnWsMessage(ev: MessageEvent) {
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
        this.counter = message.payload.counter
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
    this.ws = this.gameManagerApi.wsStateConnect();
  }

  private closeWs() {
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
  @apply container mx-auto;
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
  @apply block p-6;
  @apply border-4 rounded border-blue-300;
  width: 80px;
  height: 80px;
  text-align: center;
  transition: top 1s cubic-bezier(0, -0.18, 0.26, -0.43);
  transition-delay: 2s;
}

.cards-list__item--transition_top {
  transition: top 0.5s ease-in;
}

.cards-list__item--fail {
  @apply bg-red-700 border-red-900 text-white;
}

.cards-list__item--win {
  @apply bg-green-600 border-green-800;
}
</style>
