import { defineConfig, loadEnv, UserConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// 如果你需要使用 console.log 调试 resolve 或其他工具
console.log(resolve)

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd())

  // 返回 Vite 配置对象
  const config: UserConfig = {
    plugins: [vue()],
    server: {
      host: '0.0.0.0',
      port: 5180,
      proxy: {
        '/api/': {
          target: 'https://mock.mengxuegu.com/mock/674e5c6cc2c0134bc13ef895/traffica',
          changeOrigin: true,
          // rewrite: (path) => path.replace(/^\/api/, ''), // 可选：如果你需要重写路径
        },
        // '/gosmo/': {
        //   target: 'http://47.94.96.190:8001',
        //   changeOrigin: true,
        //   // rewrite: (path) => path.replace(/^\/api/, ''), // 可选：如果你需要重写路径
        // }
      }
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    }
  }

  return config
})
