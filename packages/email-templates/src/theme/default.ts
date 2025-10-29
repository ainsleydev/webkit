import type { EmailTheme } from './types.js';

/**
 * Default theme configuration for email templates.
 * Based on ainsley.dev design system with Search Spares colour palette.
 */
export const defaultTheme: EmailTheme = {
	colours: {
		text: {
			heading: '#1c1c1c', // $text-heading: $colour-greyscale-900
			body: '#7f7f7f', // $text-body: $colour-greyscale-500
			action: '#ff0000', // $text-action: $colour-racing-red-500
			negative: '#ffffff', // $text-negative: $colour-base-white
			darkMode: '#999999', // $text-dark-mode: $colour-greyscale-400
		},
		background: {
			white: '#ffffff', // $surface-white
			grey: '#b0b0b0', // $surface-grey / $colour-greyscale-300
			greyLight: '#f2f2f2', // $surface-grey-light / $colour-lights-500
			red: '#ff0000', // $surface-red / $colour-racing-red-500
			black: '#1c1c1c', // $surface-black / $colour-greyscale-900
		},
		border: {
			grey: '#d0d0d0', // $border-grey / $colour-greyscale-200
			black: '#1c1c1c', // $border-black / $colour-greyscale-900
			white: '#ffffff', // $border-white / $colour-base-white
			darkMode: '#333333', // $border-dark-mode / $colour-greyscale-800
		},
	},
	branding: {
		companyName: 'Search Spares',
		logoUrl: '/logo.png',
		logoWidth: 120,
		footerText: 'All rights reserved.',
		websiteUrl: 'https://ainsley.dev',
	},
};
