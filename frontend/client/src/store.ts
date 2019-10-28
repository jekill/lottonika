import Vue from 'vue';
import Vuex from 'vuex';
import {CardDto} from '@/models/CardDto';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    card: undefined,
  } as State,
  mutations: {
    setCard(state: State, value: CardDto) {
      state.card = value;
    },
  },
  actions: {
  },
});

interface State {
  card?: CardDto;
}
