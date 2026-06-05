import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), 'VITE_')

  return {
    resolve: {
      alias: {
        '@ui': resolve('../repertoire.ui/src')
      }
    },
    plugins: [
      react({
        babel: {
          plugins: [
            ['babel-plugin-react-compiler', {}]
          ]
        }
      })
    ],
    server: {
      port: parseInt(env.VITE_APPLICATION_PORT)
    }
  }
})
