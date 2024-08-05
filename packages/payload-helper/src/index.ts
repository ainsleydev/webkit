import type { Config } from 'payload';
import { fieldMapper, schemas } from './plugin/schema.js';
import type { PayloadHelperPluginConfig } from './types.js';
import env from './util/env.js';

/**
 * Payload Helper Plugin for websites at ainsley.dev
 *
 * @constructor
 * @param pluginOptions
 */
export const payloadHelper =
	(pluginOptions: PayloadHelperPluginConfig) =>
	(incomingConfig: Config): Config => {
		const genGoLang = env.bool('GEN_GOLANG', false);
		if (genGoLang) {
			incomingConfig.typescript = {
				...incomingConfig.typescript,
				schema: schemas,
			};
			// biome-ignore lint/style/noParameterAssign: Need to change field mapper.
			incomingConfig = fieldMapper(incomingConfig);
		}

		incomingConfig.typescript = incomingConfig.typescript || {};
		incomingConfig.typescript.outputFile = './src/types/payload.ts';

		return incomingConfig;
	};
