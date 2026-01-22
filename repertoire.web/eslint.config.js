import js from '@eslint/js'
import globals from 'globals'
import tslint from 'typescript-eslint'

export default tslint.config({
  files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
  ignores: ['node_modules', 'dist', '.gitignore'],
  extends: [js.configs.recommended, ...tslint.configs.recommended],
  settings: { react: { version: 'detect' } },
  languageOptions: {
    ecmaVersion: 2020,
    globals: globals.browser
  }
})
