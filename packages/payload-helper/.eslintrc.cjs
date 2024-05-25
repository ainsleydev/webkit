/** @type {import('prettier').Config} */
module.exports = {
	root: true,
	parser: '@typescript-eslint/parser',
	plugins: [
		'@typescript-eslint',
		'prettier',
		'plugin:import/errors',
		'plugin:import/warnings',
		'plugin:import/typescript',
	],
	extends: [
		'eslint:recommended',
		'plugin:@typescript-eslint/eslint-recommended',
		'plugin:@typescript-eslint/recommended',
	],
	rules: {
		// Disable some default ESLint rules that conflict with TypeScript
		semi: ['error', 'never'], // Enforce no semicolons (code style preference)
		'no-undef': 'off', // TypeScript handles undefined checks
		'no-unused-vars': 'off', // TypeScript handles unused variable checks

		// Enforce specific TypeScript ESLint rules
		'@typescript-eslint/no-explicit-any': 'error', // Disallow the 'any' type (encourages type safety)
		'@typescript-eslint/no-unused-vars': ['error'], // Enforce no unused variables in TypeScript
		'@typescript-eslint/explicit-module-boundary-types': 'off', // Explicit module boundaries (optional, can be enabled for better type checking)
	},
};
