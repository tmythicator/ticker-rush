import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '../', '')
  const exchangePort = parseInt(env.SERVER_PORT) || 8081
  const clientPort = parseInt(env.CLIENT_PORT) || 5173
  const apiUrl = `http://localhost:${exchangePort}/api`

  return {
    plugins: [react()],
    define: {
      'import.meta.env.VITE_API_URL': JSON.stringify(apiUrl),
    },
    server: {
      port: clientPort,
      host: true,
    },
  }
})
