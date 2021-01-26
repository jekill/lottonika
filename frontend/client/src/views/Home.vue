<template>
  <div class="container mx-auto">
    <div class="navbar navbar--top">
      <router-link :to="{name: 'dashboard'}">
        Dashboard
      </router-link>
    </div>
    <div class="home">
      <button
          class="button--green"
          :class="{'button--loading': isCreating}"
          @click="handleEnterClick"
      >
        <fa v-if="isCreating" icon="spinner" spin/>
        {{ $t('button.enter-game.text') }}
      </button>
    </div>
  </div>
</template>

<script lang="ts">
  import {Component, Vue, Inject} from 'vue-property-decorator';
  import {GameApi} from '@/services/GameApi';

  @Component({})
  export default class Home extends Vue {
    @Inject() public gameApi!: GameApi;

    private isCreating: boolean = false;

    public async handleEnterClick() {
      try {
        this.isCreating = true;
        const cardDto = await this.gameApi.createCard();
        this.$router.push({ name: 'card', params: { id: cardDto.id } });
      } catch (e) {
        console.error('__ERROR', e);
      } finally {
        this.isCreating = false;
      }
    }
  }
</script>
<style scoped>
  .home {
    min-height: 300px;
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

  .button--loading, .button--loading:hover {
    @apply bg-gray-300;
  }

</style>
<style>
.navbar{
  @apply bg-gray-200 px-4 py-1;
}
</style>
