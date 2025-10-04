/** @type {import("prettier").Config} */
const config = {
	useTabs: true,
	tabWidth: 4,
	singleQuote: true,
	trailingComma: 'all',
	printWidth: 100,
	overrides: [
		{
			files: '*.yaml',
			options: {
				singleQuote: false,
			},
		},
		{
			files: '*.yml',
			options: {
				singleQuote: false,
			},
		},
	],
};

export default config;
