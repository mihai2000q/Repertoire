import js from '@eslint/js'
import globals from 'globals'
import tslint from 'typescript-eslint'
import reactLint from 'eslint-plugin-react'
import prettier from 'eslint-config-prettier'
import react from "@vitejs/plugin-react";

export default tslint.config({
  files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
  ignores: ['node_modules', 'dist', '.gitignore'],
  extends: [
    js.configs.recommended,
    ...tslint.configs.recommended,
    reactLint.configs.flat.recommended,
    reactLint.configs.flat['jsx-runtime'],
    prettier
  ],
  languageOptions: {
    ...reactLint.configs.flat.recommended.languageOptions,
    ecmaVersion: 2022,
    globals: {
      ...globals.serviceworker,
      ...globals.browser
    }
  },
  plugins: { viteJs_react: react() },
  settings: { react: { version: 'detect' } }
})
