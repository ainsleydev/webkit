import type { CollectionConfig, Config } from 'payload';
import { cacheHookCollections, cacheHookGlobals } from './plugin/hooks.js';
import { injectAdminLogo } from './plugin/logo.js';
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

		// Inject admin Logo component if adminLogo config is provided
		if (pluginOptions.adminLogo) {
			config = injectAdminLogo(config, pluginOptions.adminLogo, pluginOptions.siteName);
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

export type { LogoConfig, LogoProps } from './admin/components/Logo.js';
export type { AdminLogoConfig, PayloadHelperPluginConfig } from './types.js';
