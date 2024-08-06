import type { Field } from 'payload';
import type { CollectionAfterChangeHook, GlobalAfterChangeHook } from 'payload';

/**
 * AfterChangeHook is a hook for either the Collection or Global.
 */
type AfterChangeHook = CollectionAfterChangeHook | GlobalAfterChangeHook;

/**
 * Cache hook is responsible for notifying the web server of changes
 * on Collections or Globals as defined by the endpoint.
 */
export const cacheHook = (
	endpoint: string,
	slug: string,
	fields: Field[],
	isCollection: boolean,
): AfterChangeHook => {
	//@ts-expect-error
	return async ({ doc, previousDoc, operation }) => {
		if (operation !== 'update') {
			return;
		}
		try {
			const response = await fetch(endpoint, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					slug: slug,
					fields: fields,
					type: isCollection ? 'collection' : 'global',
					doc: doc,
					prevDoc: previousDoc,
				}),
			});
			const json = await response.json();
			console.log('Webhook response', json);
		} catch (err) {
			console.error('Webhook error', err);
		}
	};
};
