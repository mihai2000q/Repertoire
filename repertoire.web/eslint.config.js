import js from "@eslint/js";
import globals from "globals";
import tslint from "typescript-eslint";

export default tslint.config({
  files: ["**/*.{ts,tsx}"],
  ignores: ["node_modules", "dist", ".gitignore", "README.md"],
  extends: [js.configs.recommended, ...tslint.configs.recommended],
  languageOptions: {
    ecmaVersion: 2020,
    globals: globals.browser,
  },
  rules: {
    "react-refresh/only-export-components": [
      "warn",
      { allowConstantExport: true },
    ],
  },
});
