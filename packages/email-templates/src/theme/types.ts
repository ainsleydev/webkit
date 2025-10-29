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
		dark: string;
		darker: string;
		highlight: string;
		accent: string;
	};
	border: {
		light: string;
		medium: string;
		dark: string;
		inverse: string;
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
	colours?: {
		text?: Partial<EmailColours['text']>;
		background?: Partial<EmailColours['background']>;
		border?: Partial<EmailColours['border']>;
	};
	branding?: Partial<EmailBranding>;
};
