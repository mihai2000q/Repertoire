import reactLint from 'eslint-plugin-react/configs/recommended'
import reactLintJsxRuntime from 'eslint-plugin-react/configs/jsx-runtime'
import prettier from 'prettier'
import electronToolkitPrettier from '@electron-toolkit/eslint-config-prettier'

export default [
  { files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'] },
  reactLint,
  reactLintJsxRuntime,
  prettier,
  electronToolkitPrettier
]
