import type { CollectionConfig, Config } from 'payload';
import { injectAdminIcon, injectAdminLogo } from './plugin/admin.js';
import { cacheHookCollections, cacheHookGlobals } from './plugin/hooks.js';
import type { PayloadHelperPluginConfig } from './types.js';

/**
 * Payload Helper Plugin for websites at ainsley.dev
 *
 * @constructor
 * @param pluginOptions
 */
export const payloadHelper =
	(pluginOptions: PayloadHelperPluginConfig) =>
	(incomingConfig: Config): Config => {
		// TODO: Validate Config

		let config = incomingConfig;

		// Update typescript generation file
		config.typescript = config.typescript || {};
		config.typescript.outputFile = './src/types/payload.ts';

		// Inject admin Logo component if logo config is provided
		if (pluginOptions.admin?.logo) {
			config = injectAdminLogo(config, pluginOptions.admin.logo, pluginOptions.siteName);
		}

		// Inject admin Icon component if icon config is provided
		if (pluginOptions.admin?.icon) {
			config = injectAdminIcon(config, pluginOptions.admin.icon, pluginOptions.siteName);
		}

		// Map collections & add hooks
		config.collections = (config.collections || []).map((collection): CollectionConfig => {
			if (collection.upload !== undefined && collection.upload !== true) {
				return collection;
			}

			const hooks = collection.hooks || {};

			// Add afterChange hook only if webServer is defined
			if (pluginOptions.webServer) {
				hooks.afterChange = [
					...(hooks.afterChange || []),
					cacheHookCollections({
						server: pluginOptions.webServer,
						slug: collection.slug,
						fields: collection.fields,
						isCollection: true,
					}),
				];
			}

			return {
				...collection,
				hooks,
			};
		});

		// Map globals & add hooks
		config.globals = (config.globals || []).map((global) => {
			const hooks = global.hooks || {};

			// Add afterChange hook only if webServer is defined
			if (pluginOptions.webServer) {
				hooks.afterChange = [
					...(hooks.afterChange || []),
					cacheHookGlobals({
						server: pluginOptions.webServer,
						slug: global.slug,
						fields: global.fields,
						isCollection: true,
					}),
				];
			}

			return {
				...global,
				hooks,
			};
		});

		return config;
	};

export type { IconProps } from './admin/components/Icon.js';
export type { LogoProps } from './admin/components/Logo.js';
export type {
	AdminConfig,
	AdminIconConfig,
	AdminLogoConfig,
	PayloadHelperPluginConfig,
} from './types.js';
