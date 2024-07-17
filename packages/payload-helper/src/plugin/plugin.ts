/**
 * Payload Plugin
 */
import type { Plugin, Config } from 'payload';

/**
 *
 * @param incomingConfig
 * @constructor
 */
export const WebKitPlugin: Plugin = (incomingConfig: Config): Config => {
	return incomingConfig;
};
