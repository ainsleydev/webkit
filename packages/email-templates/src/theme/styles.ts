import type { EmailTheme } from './types.js';

/**
 * Generates common style objects from theme configuration.
 * These can be used directly in React Email components.
 */
export function generateStyles(theme: EmailTheme) {
	return {
		heading: {
			color: theme.colours.text.heading,
			fontSize: '24px',
			lineHeight: '1.6',
			fontWeight: 'bold',
			marginTop: 0,
			marginBottom: '8px',
			letterSpacing: '-0.04em',
		},
		text: {
			fontSize: '16px',
			lineHeight: 1.5,
			color: theme.colours.text.body,
			marginTop: 0,
		},
		smallText: {
			fontSize: '14px',
			color: theme.colours.text.body,
		},
		linkText: {
			fontSize: '16px',
			lineHeight: 1.5,
			color: theme.colours.text.action,
			marginTop: 0,
		},
		button: {
			backgroundColor: theme.colours.background.dark,
			color: theme.colours.text.negative,
			padding: '10px 0',
			borderRadius: '6px',
			fontWeight: 'bold',
			textDecoration: 'none',
			fontSize: '14px',
			width: '100%',
		},
		hr: {
			borderColor: theme.colours.border.light,
			margin: '23px 0',
		},
		main: {
			backgroundColor: theme.colours.background.highlight,
			fontFamily: 'Arial, sans-serif',
		},
		container: {
			backgroundColor: theme.colours.background.white,
			margin: '60px auto',
			padding: '30px',
			borderRadius: '10px',
			maxWidth: '600px',
		},
		logoSection: {
			textAlign: 'center' as const,
			marginBottom: '20px',
		},
		footerText: {
			fontSize: '13px',
			color: theme.colours.text.body,
			margin: '0',
		},
	};
}
