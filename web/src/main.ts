import { createApp, type Component } from 'vue'
import { createPinia } from 'pinia'
import { plugin } from '@formkit/vue'
import PrimeVue from 'primevue/config'
import formKitConfig from '../formkit.config'
import Lara from '@primeuix/themes/lara'
import ToastService from 'primevue/toastservice'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import Tooltip from 'primevue/tooltip'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import './main.css'

import ConfirmationService from 'primevue/confirmationservice'

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(plugin, formKitConfig)
app.use(router)
app.use(i18n)
app.use(ConfirmationService)
app.use(PrimeVue, {
  theme: {
    preset: Lara,
    unstyled: true,
    options: {
      darkModeSelector: '.app-dark'
    }
  }
})
app.use(ToastService)
app.directive('tooltip', Tooltip)

app.config.errorHandler = function (err, vm, info) {
  console.log(err)
}
app.mount('#app')
