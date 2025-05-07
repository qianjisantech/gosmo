import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())
  const API_BASE_URL = env.VITE_API_BASE_URL || 'https://mock.mengxuegu.com/mock/681b2cb1d5b98b579eb29a0f/gosmo'
  const API_GOSMO_URL = env.VITE_API_GOSMO_URL || 'http://47.94.96.190:9200'
  console.log('API_BASE_URL:', API_BASE_URL) // 确认输出正确
  console.log('API_GOSMO_URL:', API_GOSMO_URL) // 确认输出正确
  return {
    plugins: [vue()],
    server: {
      host: '0.0.0.0',
      port: 5180,
      proxy: {
        '/api': {
          target: API_BASE_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, '') // 去掉 /api 前缀（根据后端要求调整）
        }
        // ,
        // '/gosmo': {
        //   target: API_GOSMO_URL,
        //   changeOrigin: true,
        //   rewrite: (path) => path.replace(/^\/gosmo/, '') // 去掉 /api 前缀（根据后端要求调整）
        // }
      }
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    }
  }
})