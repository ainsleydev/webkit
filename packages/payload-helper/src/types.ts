// import type { SEOPluginConfig } from "@payloadcms/plugin-seo/types";
import type { GlobalConfig, Tab } from 'payload';

export type SettingsConfig = {
	additionalTabs?: Tab[];
	override: (args: {
		config: GlobalConfig;
	}) => GlobalConfig;
};

export type WebServerConfig = {
	apiKey?: string;
	cacheEndpoint?: string;
};

export type PayloadHelperPluginConfig = {
	settings?: SettingsConfig;
	// seo?: (args: {
	// 	config: SEOPluginConfig;
	// }) => SEOPluginConfig;

	webServer?: WebServerConfig;
};
