import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import axios from 'axios'

import App from './App.vue'
import router from './router'

// 创建应用实例
const app = createApp(App)

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 配置axios默认值
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
axios.defaults.timeout = 10000
axios.defaults.headers.common['Content-Type'] = 'application/json'

// 使用插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus, { size: 'default' })

// 挂载应用
app.mount('#app') 