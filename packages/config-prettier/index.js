/**
 * Prettier Configuration for ainsley.dev
 *
 * @see https://prettier.io/docs/configuration
 * @see https://prettier.io/docs/sharing-configurations
 *
 * @type {import("prettier").Config}
 */
const config = {
	useTabs: true,
	singleQuote: true,
	trailingComma: 'all',
	printWidth: 100,
	tabWidth: 4,
	semi: false,
	plugins: ['prettier-plugin-svelte'],
	overrides: [
		{
			files: ['*.yml', '*.yaml'],
			options: {
				useTabs: false,
				tabWidth: 2,
			},
		},
	],
};

module.exports = config;
