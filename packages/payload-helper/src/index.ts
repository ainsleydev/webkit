import type { CollectionConfig, Config } from 'payload';
import { cacheHookCollections, cacheHookGlobals } from './plugin/hooks.js';
import { fieldMapper, schemas } from './plugin/schema.js';
import type { PayloadHelperPluginConfig } from './types.js';
import env from './util/env.js';

// export const test = (pluginOptions: PayloadHelperPluginConfig): Plugin[] => {
// 	return [
// 		seoPlugin({
// 			collections: pluginOptions?.seo?.collections,
// 			globals: pluginOptions?.seo?.globals,
// 			fields: [...SEOFields, pluginOptions.seo?.fields],
// 			tabbedUI: true,
// 			uploadsCollection: 'media',
// 			generateTitle: pluginOptions?.seo?.generateTitle ?
// 				pluginOptions?.seo?.generateTitle :
// 				({ doc }) => `${pluginOptions.siteName} — ${doc?.title?.value ?? ''}`,
// 			generateDescription: pluginOptions?.seo?.generateDescription ?
// 				pluginOptions?.seo?.generateDescription :
// 				({ doc }) => doc?.excerpt?.value,
// 		}),
// 	];
// };

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

		// Update typescript generation file
		incomingConfig.typescript = incomingConfig.typescript || {};
		incomingConfig.typescript.outputFile = './src/types/payload.ts';

		// Map collections & add hooks
		incomingConfig.collections = (incomingConfig.collections || []).map(
			(collection): CollectionConfig => {
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
			},
		);

		// Map globals & add hooks
		incomingConfig.globals = (incomingConfig.globals || []).map((global) => {
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

		return incomingConfig;
	};
