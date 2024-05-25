import globals from 'globals';
import pluginJs from '@eslint/js';
import tseslint from 'typescript-eslint';

export default [
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
			},
		},
		ignores: [
			'node_modules/*',
			'templates/'
		]
	},
	pluginJs.configs.recommended,
	...tseslint.configs.recommended,
];
