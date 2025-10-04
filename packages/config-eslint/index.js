import { fileURLToPath } from 'node:url';
import { includeIgnoreFile } from '@eslint/compat';
import js from '@eslint/js';
import prettier from 'eslint-config-prettier';
import importPlugin from 'eslint-plugin-import';
import perfectionist from 'eslint-plugin-perfectionist';
import svelte from 'eslint-plugin-svelte';
import ts from 'typescript-eslint';

/** @typedef {import('eslint').Linter.Config} FlatConfig */

/**
 * Helper to include .gitignore from project root
 */
export function withGitignore(projectRoot) {
	const gitignorePath = fileURLToPath(new URL('./.gitignore', projectRoot));
	return includeIgnoreFile(gitignorePath);
}

const baseRules = {
	'class-methods-use-this': 'off',
	curly: ['warn', 'all'],
	'arrow-body-style': 'off',
	'no-restricted-exports': ['warn', { restrictDefaultExports: { direct: true } }],
	'no-console': 'warn',
	'no-sparse-arrays': 'off',
	'no-underscore-dangle': 'off',
	'no-use-before-define': 'off',
	'object-shorthand': 'warn',
	'no-useless-escape': 'warn',
	'perfectionist/sort-objects': [
		'error',
		{
			type: 'natural',
			order: 'asc',
			partitionByComment: true,
			partitionByNewLine: true,
			groups: ['top', 'unknown'],
			customGroups: {
				top: ['_id', 'id', 'name', 'slug', 'type'],
			},
		},
	],
};

const typescriptRules = {
	'@typescript-eslint/no-explicit-any': 'warn',
	'@typescript-eslint/no-empty-object-type': 'warn',
};

/**
 * Base config (applied to all projects)
 * @type {FlatConfig[]}
 */
const baseConfig = [
	js.configs.recommended,
	...ts.configs.recommended,
	perfectionist.configs['recommended-natural'],
	prettier,
	{
		plugins: {
			import: importPlugin,
		},
		rules: {
			...baseRules,
			...typescriptRules,
		},
	},
];

/**
 * Svelte-specific config (only applied to Svelte files)
 * @type {FlatConfig[]}
 */
const svelteConfig = [
	...svelte.configs.recommended,
	...svelte.configs.prettier,
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: ts.parser,
				extraFileExtensions: ['.svelte'],
			},
		},
		rules: {
			'svelte/valid-compile': ['error', { ignoreWarnings: false }],
			'svelte/no-navigation-without-resolve': ['warn'],
			'svelte/no-at-html-tags': 'off',
		},
	},
];

/**
 * Import ordering config (applies to all file types)
 * @type {FlatConfig[]}
 */
const importConfig = [
	{
		files: ['**/*.ts', '**/*.tsx', '**/*.js', '**/*.jsx', '**/*.svelte'],
		rules: {
			'import/order': [
				'error',
				{
					groups: ['builtin', 'external', 'internal', ['parent', 'sibling', 'index']],
					'newlines-between': 'always',
					alphabetize: { order: 'asc', caseInsensitive: true },
					pathGroups: [
						{ pattern: '$lib/**', group: 'internal' },
						{ pattern: '$app/**', group: 'internal' },
					],
					pathGroupsExcludedImportTypes: ['builtin'],
				},
			],
		},
	},
];

/**
 * Complete ainsley.dev ESLint configuration
 * @type {FlatConfig[]}
 */
export default [...baseConfig, ...svelteConfig, ...importConfig];
