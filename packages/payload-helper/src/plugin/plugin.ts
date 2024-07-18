import type { Config } from 'payload';
import { fieldMapper, schemas } from './schema';
import env from '../util/env';

/**
 * Plugin Options
 */
export interface PluginOptions {
	SEOFields?: boolean;
}

/**
 * Payload Helper Plugin for websites at ainsley.dev
 *
 * @constructor
 * @param pluginOptions
 */
export const payloadHelper =
	(pluginOptions: {}) =>
	(incomingConfig: Config): Config => {
		console.log(pluginOptions);

		const genGoLang = env.bool('GEN_GOLANG', false);
		if (genGoLang) {
			incomingConfig.typescript = {
				...incomingConfig.typescript,
				schema: schemas,
			};
			incomingConfig = fieldMapper(incomingConfig);
		}

		incomingConfig.typescript = incomingConfig.typescript || {};
		incomingConfig.typescript.outputFile = './src/types/payload.ts';

		return incomingConfig;
	};
