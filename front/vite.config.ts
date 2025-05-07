import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  // ✅ 生产环境直接使用真实 API，开发环境走代理
  const API_BASE_URL = env.VITE_API_BASE_URL || 'https://mock.mengxuegu.com/mock/681b2cb1d5b98b579eb29a0f/gosmo'

  return {
    plugins: [vue()],
    server: {
      host: '0.0.0.0',
      port: 5180,
      proxy: {
        '/api/': {
          target: 'https://mock.mengxuegu.com/mock/681b2cb1d5b98b579eb29a0f/gosmo',
          changeOrigin: true,
          // rewrite: (path) => path.replace(/^\/api/, ''),
        }
      }
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    },
    define: {
      'process.env.VITE_API_BASE_URL': JSON.stringify(API_BASE_URL)
    }
  }
})