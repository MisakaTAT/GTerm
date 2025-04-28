import '@/styles/basic.less';
import { LogInfo } from '@wailsApp/runtime/runtime';
import { createPinia } from 'pinia';
import { createApp } from 'vue';
import router from '@/router';
import { i18n } from '@/utils/i18n';
import App from './App.vue';

async function setupApp() {
  const app = createApp(App);
  app.use(router);
  app.use(i18n);
  app.use(createPinia());
  app.mount('#app');
}

setupApp().then(() => LogInfo('App is ready'));
