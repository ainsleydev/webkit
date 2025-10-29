import type { EmailTheme } from './types.js';

/**
 * Default theme configuration for email templates.
 * Based on ainsley.dev design system with Search Spares colour palette.
 */
export const defaultTheme: EmailTheme = {
	// Taken from: https://github.com/ainsleydev/website/blob/main/assets/scss/abstracts/_variables.scss
	colours: {
		text: {
			heading: '#0a0a0a', // $black / --colour-black
			body: '#9a9a9a', // $copy-dark-bg / --colour-paragraph
			action: '#ff5043', // $orange / --colour-orange
			negative: '#ffffff', // $white / --colour-white
			darkMode: '#595959', // $copy-light-bg / --colour-copy-light
		},
		background: {
			white: '#ffffff', // $white / --colour-white
			dark: '#0a0a0a', // $black / --colour-black
			darker: '#0f0f0f', // $grey-dark / --colour-grey-dark
			highlight: '#171717', // $grey-light / --colour-grey-light
			accent: '#ff5043', // $orange / --colour-orange
		},
		border: {
			light: '#2b2b2b', // --table-border-colour
			medium: 'rgba(10, 10, 10, 0.15)', // rgba($black, $alpha-standard)
			dark: '#0a0a0a', // $black / --colour-black
			inverse: '#ffffff', // $white / --colour-white
		},
	},
	branding: {
		companyName: 'ainsley.dev',
		logoUrl: '/logo.png',
		logoWidth: 120,
		footerText: 'All rights reserved.',
		websiteUrl: 'https://ainsley.dev',
	},
};
