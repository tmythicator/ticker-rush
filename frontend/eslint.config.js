import js from '@eslint/js';
import eslintConfigPrettier from 'eslint-config-prettier';
import cssModules from 'eslint-plugin-css-modules';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import { defineConfig } from 'eslint/config';
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default defineConfig([
    {
        ignores: ['dist', 'src/lib/proto', 'vite.config.ts', 'eslint.config.js', 'stylelint.config.js'],
    },
    {
        files: ['**/*.{ts,tsx}'],
        extends: [
            js.configs.recommended,
            ...tseslint.configs.recommended,
            reactHooks.configs.flat.recommended,
            reactRefresh.configs.vite,
            jsxA11y.flatConfigs.recommended,
        ],
        plugins: {
            'css-modules': cssModules,
        },
        languageOptions: {
            ecmaVersion: 2020,
            globals: globals.browser,
        },

        rules: {
            'css-modules/no-unused-class': 'error',
            'css-modules/no-undef-class': 'error',

            '@typescript-eslint/no-unused-vars': [
                'error',
                {
                    argsIgnorePattern: '^_',
                    varsIgnorePattern: '^_',
                },
            ],

            '@typescript-eslint/consistent-type-imports': [
                'error',
                {
                    prefer: 'type-imports',
                },
            ],

            complexity: ['error', { max: 10 }],
            'no-debugger': 'error',
        },
    },
    {
        files: ['**/*.test.{ts,tsx}', '**/*.spec.{ts,tsx}', 'src/test/**/*'],
        rules: {
            'complexity': 'off',
        },
    },
    eslintConfigPrettier,
]);
