// import type { SEOPluginConfig } from "@payloadcms/plugin-seo/types";
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

export type PayloadHelperPluginConfig = {
	siteName: string;
	settings?: SettingsConfig;
	// seo?: SEOConfig;
	webServer?: WebServerConfig;
};
