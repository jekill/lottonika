import Vue from 'vue';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faSpinner } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
// @ts-ignore
import Qui from '@qvant/qui';
import '@qvant/qui/dist/qui.css';

import App from './App.vue';
import router from './router';
import store from './store';
import i18n from './i18n';

require('@/assets/main.css');

Vue.config.productionTip = false;

Vue.component('fa', FontAwesomeIcon);
library.add(faSpinner);

Vue.use(Qui);

new Vue({
  router,
  store,
  i18n,
  render: (h) => h(App),
}).$mount('#app');
