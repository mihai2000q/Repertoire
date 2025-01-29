import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'happy-dom',
    globals: true,
    setupFiles: './src/setupTests.ts',
    unstubEnvs: true,
    testTimeout: 10_000
  }
})
