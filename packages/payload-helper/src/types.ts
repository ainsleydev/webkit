// import type { SEOPluginConfig } from "@payloadcms/plugin-seo/types";
import type { PartialEmailTheme } from '@ainsleydev/email-templates';
import type { CollectionConfig, Config, GlobalConfig, Tab, TextField, TextareaField, UploadField } from 'payload';

/**
 * Arguments passed to custom URL generator callbacks.
 */
export type EmailUrlCallbackArgs = {
	/**
	 * The authentication token for the email action.
	 */
	token: string;
	/**
	 * The Payload configuration object.
	 */
	config: Config;
	/**
	 * The collection configuration for the user being emailed.
	 */
	collection: CollectionConfig;
};

/**
 * Callback function type for generating custom email URLs.
 */
export type EmailUrlCallback = (args: EmailUrlCallbackArgs) => string;
// import type {SEOPluginConfig} from "@payloadcms/plugin-seo/dist/types.js";

/**
 * Configuration for the Settings global.
 */
export type SettingsConfig = {
	/**
	 * Additional tabs to add to the Settings global.
	 */
	additionalTabs?: Tab[];
	/**
	 * Function to override the entire Settings global configuration.
	 */
	override: (args: {
		config: GlobalConfig;
	}) => GlobalConfig;
};

/**
 * Configuration for web server cache invalidation.
 */
export type WebServerConfig = {
	/**
	 * Optional API key for authenticating cache invalidation requests.
	 */
	apiKey?: string;
	/**
	 * Base URL of the web server.
	 */
	baseURL: string;
	/**
	 * Endpoint path for cache invalidation.
	 */
	cacheEndpoint: string;
};

// export type SEOConfig = Omit<SEOPluginConfig, 'uploadsCollection' | 'tabbedUI'>;
//
// export type S3Config = {
// 	enabled: boolean;
// 	bucket: string
// 	config: AWS.S3ClientConfig;
// }

/**
 * Configuration for the admin panel logo.
 */
export type AdminLogoConfig = {
	/**
	 * Path to the logo image file.
	 */
	path: string;
	/**
	 * Optional path to the dark mode logo image file.
	 */
	darkModePath?: string;
	/**
	 * Optional width of the logo in pixels.
	 */
	width?: number;
	/**
	 * Optional height of the logo in pixels.
	 */
	height?: number;
	/**
	 * Optional alt text for the logo image.
	 */
	alt?: string;
	/**
	 * Optional CSS class name for the logo.
	 */
	className?: string;
};

/**
 * Configuration for the admin panel icon/favicon.
 */
export type AdminIconConfig = {
	/**
	 * Path to the icon image file.
	 */
	path: string;
	/**
	 * Optional path to the dark mode icon image file.
	 */
	darkModePath?: string;
	/**
	 * Optional width of the icon in pixels.
	 */
	width?: number;
	/**
	 * Optional height of the icon in pixels.
	 */
	height?: number;
	/**
	 * Optional alt text for the icon image.
	 */
	alt?: string;
	/**
	 * Optional CSS class name for the icon.
	 */
	className?: string;
};

/**
 * Configuration for admin panel customization.
 */
export type AdminConfig = {
	/**
	 * Optional logo configuration for the admin panel.
	 */
	logo?: AdminLogoConfig;
	/**
	 * Optional icon/favicon configuration for the admin panel.
	 */
	icon?: AdminIconConfig;
};

/**
 * Content overrides for customizing email template text.
 */
export type EmailContentOverrides = {
	/**
	 * Optional preview text shown in email clients before opening.
	 */
	previewText?: string;
	/**
	 * Optional heading text displayed at the top of the email.
	 */
	heading?: string;
	/**
	 * Optional body text content of the email.
	 */
	bodyText?: string;
	/**
	 * Optional button text for the call-to-action button.
	 */
	buttonText?: string;
};

/**
 * Configuration for email templates used in authentication flows.
 */
export type EmailConfig = {
	/**
	 * Optional front-end URL override for email links. If not provided, uses Payload's serverUrl.
	 */
	frontEndUrl?: string;

	/**
	 * Optional theme customization for email templates (colors, branding, etc.).
	 */
	theme?: PartialEmailTheme;

	/**
	 * Optional content overrides for the forgot password email template.
	 */
	forgotPassword?: EmailContentOverrides;

	/**
	 * Optional content overrides for the verify account email template.
	 */
	verifyAccount?: EmailContentOverrides;

	/**
	 * Optional callback to generate a custom forgot password URL.
	 * When provided, this overrides the default URL generation.
	 *
	 * @example
	 * ```ts
	 * forgotPasswordUrl: ({ token, config, collection }) =>
	 *   `https://myapp.com/auth/reset-password?token=${token}`
	 * ```
	 */
	forgotPasswordUrl?: EmailUrlCallback;

	/**
	 * Optional callback to generate a custom verify account URL.
	 * When provided, this overrides the default URL generation.
	 *
	 * @example
	 * ```ts
	 * verifyAccountUrl: ({ token, config, collection }) =>
	 *   `https://myapp.com/auth/verify?token=${token}&collection=${collection.slug}`
	 * ```
	 */
	verifyAccountUrl?: EmailUrlCallback;
};

/**
 * Main configuration object for the Payload Helper plugin.
 */
export type PayloadHelperPluginConfig = {
	/**
	 * The name of the site, used throughout the admin panel and emails.
	 */
	siteName: string;
	/**
	 * Optional settings global configuration.
	 */
	settings?: SettingsConfig;
	// seo?: SEOConfig;
	/**
	 * Optional web server configuration for cache invalidation.
	 */
	webServer?: WebServerConfig;
	/**
	 * Optional admin panel customization (logo, icon).
	 */
	admin?: AdminConfig;
	/**
	 * Optional email template configuration for authentication emails.
	 */
	email?: EmailConfig;
};
