import Vue from 'vue';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faSpinner } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import App from './App.vue';
import router from './router';
import store from './store';
import i18n from './i18n';

require('@/assets/main.css');

Vue.config.productionTip = false;

Vue.component('fa', FontAwesomeIcon);
library.add(faSpinner);

new Vue({
  router,
  store,
  i18n,
  render: (h) => h(App),
}).$mount('#app');
