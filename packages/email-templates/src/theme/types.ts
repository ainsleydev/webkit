/**
 * Colour configuration for email templates.
 */
export interface EmailColours {
	text: {
		heading: string;
		body: string;
		action: string;
		negative: string;
		darkMode?: string;
	};
	background: {
		white: string;
		grey: string;
		greyLight: string;
		red: string;
		black: string;
	};
	border: {
		grey: string;
		black: string;
		white: string;
		darkMode?: string;
	};
}

/**
 * Branding configuration for email templates.
 */
export interface EmailBranding {
	companyName: string;
	logoUrl: string;
	logoWidth: number;
	footerText?: string;
	websiteUrl?: string;
}

/**
 * Complete theme configuration for email templates.
 * All properties support partial overrides which will be merged with defaults.
 */
export interface EmailTheme {
	colours: EmailColours;
	branding: EmailBranding;
}

/**
 * Partial theme configuration allowing selective overrides.
 */
export type PartialEmailTheme = {
	colours?: Partial<EmailColours> & {
		text?: Partial<EmailColours['text']>;
		background?: Partial<EmailColours['background']>;
		border?: Partial<EmailColours['border']>;
	};
	branding?: Partial<EmailBranding>;
};
