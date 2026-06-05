import globals from 'globals'
import jsLint from '@eslint/js'
import tsLint from 'typescript-eslint'
import reactLint from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'
import prettier from 'eslint-config-prettier'

export default tsLint.config(
  {
    files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
    ignores: ['node_modules', '.gitignore'],
    languageOptions: {
      globals: {
        ...globals.serviceworker,
        ...globals.browser,
        ...globals.node
      }
    },
    settings: { react: { version: 'detect' } }
  },
  jsLint.configs.recommended,
  ...tsLint.configs.recommended,
  reactLint.configs.flat.recommended,
  reactLint.configs.flat['jsx-runtime'],
  reactHooks.configs.flat.recommended,
  {
    rules: {
      'react/react-in-jsx-scope': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'no-console': 'warn',
      '@typescript-eslint/no-require-imports': 'warn',
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'off',
      'react-hooks/set-state-in-effect': 'off',
      'react-hooks/refs': 'off'
    }
  },
  prettier
)
