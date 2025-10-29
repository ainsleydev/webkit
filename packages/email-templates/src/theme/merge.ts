import type { EmailTheme, PartialEmailTheme } from './types.js';
import { defaultTheme } from './default.js';

/**
 * Deep merges a partial theme with the default theme.
 * Allows for selective overrides of theme properties.
 *
 * @param partial - Partial theme configuration to merge with defaults
 * @returns Complete theme with overrides applied
 */
export function mergeTheme(partial?: PartialEmailTheme): EmailTheme {
	if (!partial) {
		return defaultTheme;
	}

	return {
		colours: {
			text: {
				...defaultTheme.colours.text,
				...partial.colours?.text,
			},
			background: {
				...defaultTheme.colours.background,
				...partial.colours?.background,
			},
			border: {
				...defaultTheme.colours.border,
				...partial.colours?.border,
			},
		},
		branding: {
			...defaultTheme.branding,
			...partial.branding,
		},
	};
}
