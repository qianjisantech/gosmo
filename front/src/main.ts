import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import Antd from 'ant-design-vue';
import VChart from 'vue-echarts';
import 'echarts';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 创建 Vue 应用
const app = createApp(App);

// 注册插件和组件
app.use(router); // 使用路由
app.use(Antd); // 使用 Ant Design Vue
app.use(ElementPlus)
app.component('VChart', VChart); // 注册 VChart 组件

// 挂载应用
app.mount('#app');
