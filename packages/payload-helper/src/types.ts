// import type { SEOPluginConfig } from "@payloadcms/plugin-seo/types";
import type { PartialEmailTheme } from '@ainsleydev/email-templates';
import type { GlobalConfig, Tab, TextField, TextareaField, UploadField } from 'payload';
// import type {SEOPluginConfig} from "@payloadcms/plugin-seo/dist/types.js";

export type SettingsConfig = {
	additionalTabs?: Tab[];
	override: (args: {
		config: GlobalConfig;
	}) => GlobalConfig;
};

export type WebServerConfig = {
	apiKey?: string;
	baseURL: string;
	cacheEndpoint: string;
};

// export type SEOConfig = Omit<SEOPluginConfig, 'uploadsCollection' | 'tabbedUI'>;
//
// export type S3Config = {
// 	enabled: boolean;
// 	bucket: string
// 	config: AWS.S3ClientConfig;
// }

export type AdminLogoConfig = {
	path: string;
	darkModePath?: string;
	width?: number;
	height?: number;
	alt?: string;
	className?: string;
};

export type AdminIconConfig = {
	path: string;
	darkModePath?: string;
	width?: number;
	height?: number;
	alt?: string;
	className?: string;
};

export type AdminConfig = {
	logo?: AdminLogoConfig;
	icon?: AdminIconConfig;
};

export type EmailContentOverrides = {
	previewText?: string;
	heading?: string;
	bodyText?: string;
	buttonText?: string;
};

export type EmailConfig = {
	/**
	 * Optional front-end URL override. If not provided, uses Payload's serverUrl.
	 */
	frontEndUrl?: string;

	/**
	 * Optional theme customization for email templates.
	 */
	theme?: PartialEmailTheme;

	/**
	 * Optional content overrides for the forgot password email.
	 */
	forgotPassword?: EmailContentOverrides;

	/**
	 * Optional content overrides for the verify account email.
	 */
	verifyAccount?: EmailContentOverrides;
};

export type PayloadHelperPluginConfig = {
	siteName: string;
	settings?: SettingsConfig;
	// seo?: SEOConfig;
	webServer?: WebServerConfig;
	admin?: AdminConfig;
	email?: EmailConfig;
};
