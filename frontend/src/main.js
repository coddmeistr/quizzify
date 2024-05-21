import { createApp } from 'vue'

import store from './store'

// Vuetify
import 'vuetify/styles'
import vuetify from './plugins/vuetify'

// Components
import App from './App.vue'
import router from './router'


createApp(App).use(router).use(vuetify).use(store).mount('#app')

