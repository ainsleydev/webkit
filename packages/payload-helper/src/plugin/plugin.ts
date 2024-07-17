import type { Config } from 'payload';
import env from '../util/env';
import { fieldMapper, schemas } from './schema';

/**
 * Plugin Options
 */
export interface PluginOptions {
	SEOFields: boolean;
}

/**
 * Payload Helper Plugin for websites at ainsley.dev
 *
 * @constructor
 * @param pluginOptions
 */
export const payloadHelper =
	(pluginOptions: PluginOptions) =>
	(incomingConfig: Config): Config => {
		console.log(pluginOptions);

		const genGoLang = env.bool('GEN_GOLANG', false);
		if (genGoLang) {
			incomingConfig.typescript = {
				...incomingConfig.typescript,
				schema: schemas(incomingConfig),
			};
			incomingConfig = fieldMapper(incomingConfig);
		}

		if (!incomingConfig.typescript || incomingConfig.typescript.outputFile === undefined) {
			incomingConfig.typescript.outputFile = './types/payload.ts';
		}

		return incomingConfig;
	};
