import { createApp } from 'vue'
import { createPinia } from 'pinia'

import "go-captcha-vue/dist/style.css"
import GoCaptcha from "go-captcha-vue"

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(GoCaptcha)

app.mount('#app')
