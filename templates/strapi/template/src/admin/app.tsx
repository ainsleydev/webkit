import favicon from './extensions/favicon.ico';
import logo from './extensions/logo.svg';
import symbol from './extensions/symbol.svg';

/**
 * Customisation of admin panel.
 *
 * TODO Add notes on faviocon etc
 *
 * @see: https://fffuel.co/cccolor/
 * @see: https://docs.strapi.io/dev-docs/admin-panel-customization
 * @see: https://github.com/strapi/design-system/blob/main/packages/strapi-design-system/src/themes/
 */
export default {
	config: {
		locales: [],
		head: {
			//favicon: favicon,
		},
		auth: {
			//logo: symbol,
		},
		menu: {
			//logo: symbol,
		},
		tutorials: false,
		theme: {
			light: {
				// NOTE: 100s include the background colour, 500+ is the text.
				colors: {},
			},
		},
	},
	bootstrap(app: unknown) {
		console.log(app);
	},
};
