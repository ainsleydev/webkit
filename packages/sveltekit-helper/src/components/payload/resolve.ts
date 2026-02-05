/**
 * Resolves a Payload relationship field value to its populated object.
 * Returns null if the value is an unresolved ID (number), null or undefined.
 *
 * @param {number | null | undefined | T} item - The relationship field value to resolve.
 * @returns {T | null} The populated object, or null if unresolved.
 */
export const resolveItems = <T>(item: number | null | undefined | T): T | null => {
	if (item === null || item === undefined || typeof item === 'number') return null;
	return item;
};
