import reactLint from 'eslint-plugin-react'
import jsLint from '@eslint/js'
import tsLint from 'typescript-eslint'
import prettier from 'eslint-config-prettier'
import electronToolkitPrettier from '@electron-toolkit/eslint-config-prettier'

export default tsLint.config(
  { files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'] },
  { ignores: ['node_modules', 'dist', 'out', '.gitignore'] },
  { settings: { react: { version: 'detect' } } },
  reactLint.configs.flat.recommended,
  jsLint.configs.recommended,
  ...tsLint.configs.recommended,
  electronToolkitPrettier,
  prettier
)
