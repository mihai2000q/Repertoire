import globals from 'globals'
import jsLint from '@eslint/js'
import tsLint from 'typescript-eslint'
import reactLint from 'eslint-plugin-react'
import prettier from 'eslint-config-prettier'
import react from '@vitejs/plugin-react'

export default tsLint.config(
  {
    files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
    extends: [
      jsLint.configs.recommended,
      ...tsLint.configs.recommended,
      reactLint.configs.flat.recommended,
      reactLint.configs.flat['jsx-runtime'],
      prettier
    ],
    ignores: ['node_modules', '.gitignore'],
    languageOptions: {
      ...reactLint.configs.flat.recommended.languageOptions,
      globals: {
        ...globals.serviceworker,
        ...globals.browser,
        ...globals.node
      }
    },
    plugins: { viteJs_react: react() },
    settings: { react: { version: 'detect' } },
    rules: {
      'react/react-in-jsx-scope': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'no-console': 'warn',
      '@typescript-eslint/no-require-imports': 'warn'
    }
  }
)
