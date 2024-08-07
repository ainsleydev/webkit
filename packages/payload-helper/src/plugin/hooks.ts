import type { Field, Payload } from 'payload';
import type { CollectionAfterChangeHook, GlobalAfterChangeHook } from 'payload';
import type { WebServerConfig } from '../types.js';

/**
 * TODO
 */
export type CacheBustConfig = {
	server?: WebServerConfig;
	slug: string;
	fields: Field[];
	isCollection: boolean;
};

/**
 * Cache hook is responsible for notifying the web server of changes
 * on Collections or Globals as defined by the endpoint.
 */
const cacheBust = async (
	config: CacheBustConfig,
	payload: Payload,
	doc: unknown,
	previousDoc: unknown,
) => {
	const logger = payload.logger;

	const endpoint =
		new URL(config?.server?.cacheEndpoint ?? '', config?.server?.baseURL ?? '').href ?? '';

	try {
		const response = await fetch(endpoint, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				slug: config.slug,
				fields: config.fields,
				type: config.isCollection ? 'collection' : 'global',
				doc: doc,
				prevDoc: previousDoc,
			}),
		});
		logger.info(`Webhook response status: ${response.status}`);
	} catch (err) {
		logger.error(`Webhook error ${err}`);
	}
};

export const cacheHookCollections = (config: CacheBustConfig): CollectionAfterChangeHook => {
	return async ({ req, doc, previousDoc, operation }) => {
		if (operation !== 'update') {
			return;
		}
		await cacheBust(config, req.payload, doc, previousDoc);
	};
};

export const cacheHookGlobals = (config: CacheBustConfig): GlobalAfterChangeHook => {
	return async ({ req, doc, previousDoc }) => {
		await cacheBust(config, req.payload, doc, previousDoc);
	};
};
