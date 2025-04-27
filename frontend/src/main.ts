import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import i18n from './i18n'
import { setupElementPlus } from './plugins/element-plus'
import App from './App.vue'
import './style.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

// 设置Element Plus
setupElementPlus(app)

app.mount('#app')