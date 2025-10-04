import js from '@eslint/js';
import payloadPlugin from '@payloadcms/eslint-plugin';
import prettier from 'eslint-config-prettier';
import perfectionist from 'eslint-plugin-perfectionist';
import svelte from 'eslint-plugin-svelte';
import { globalIgnores } from 'eslint/config';
import globals from 'globals';
import ts from 'typescript-eslint';

/**
 * Base ainsley.dev configuration for ESLint, it can be
 * used in any project using JS/TS.
 */
export const baseConfig = [
	js.configs.recommended,
	...ts.configs.recommended,

	{
		plugins: {
			perfectionist,
		},
		languageOptions: {
			globals: { ...globals.browser, ...globals.node },
		},
		rules: {
			// TypeScript rules
			'@typescript-eslint/no-empty-object-type': 'warn',
			'@typescript-eslint/no-explicit-any': 'off',

			// ts-expect preferred over ts-ignore. It will error if the expected error is no longer present.
			'@typescript-eslint/ban-ts-comment': 'warn',

			// By default, it errors for unused variables. This is annoying, warnings are enough.
			'@typescript-eslint/no-unused-vars': [
				'warn',
				{
					vars: 'all',
					args: 'after-used',
					ignoreRestSiblings: false,
					argsIgnorePattern: '^_',
					varsIgnorePattern: '^_',
					destructuredArrayIgnorePattern: '^_',
					caughtErrorsIgnorePattern: '^(_|ignore)',
				},
			],

			// Disable no-undef for TypeScript files
			'no-undef': 'off',

			// Perfectionist import sorting
			'perfectionist/sort-imports': [
				'error',
				{
					type: 'alphabetical',
					order: 'asc',
					ignoreCase: true,
					newlinesBetween: 'always',
					internalPattern: ['^@/.*', '^\\$lib/.*'],
					groups: [
						'builtin',
						'external',
						'internal',
						['parent', 'sibling', 'index'],
						'object',
						'type',
					],
				},
			],
		},
	},

	// Prettier config (should be last to override formatting rules)
	prettier,

	// Common ignore patterns that apply to most projects
	globalIgnores([
		'**/node_modules/**',
		'**/dist/**',
		'**/build/**',
		'**/.next/**',
		'**/.svelte-kit/**',
		'**/coverage/**',
		'**/.turbo/**',
		'**/migrations/**',
		'**/importMap.js',
	]),
];

/**
 * Svelte Configuration
 */
export const svelteConfig = [
	...svelte.configs.recommended,
	...svelte.configs.prettier,

	{
		files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
		languageOptions: {
			parserOptions: {
				projectService: true,
				extraFileExtensions: ['.svelte'],
				parser: ts.parser,
			},
		},
		rules: {
			'svelte/valid-compile': ['error', { ignoreWarnings: false }],
			'svelte/no-navigation-without-resolve': ['off'],
			'svelte/no-at-html-tags': 'off',
		},
	},
];

/**
 * Payload CMS configuration.
 */
export const payloadConfig = [
	{
		plugins: {
			payload: payloadPlugin,
		},
	},
];

export default baseConfig;
