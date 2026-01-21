import globals from 'globals'
import jsLint from '@eslint/js'
import tsLint from 'typescript-eslint'
import reactLint from 'eslint-plugin-react'
import prettier from 'eslint-config-prettier'

export default tsLint.config(
  jsLint.configs.recommended,
  ...tsLint.configs.recommended,
  reactLint.configs.flat.recommended,
  {
    files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
    ignores: ['node_modules', '.gitignore', 'README.md'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node
      }
    },
    rules: {
      'react/react-in-jsx-scope': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'no-console': 'warn',
      '@typescript-eslint/no-require-imports': 'warn'
    }
  },
  prettier
)
