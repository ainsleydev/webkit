/**
 * Serialises a value as a JSON-LD script tag for embedding
 * structured data in HTML.
 *
 * @param {unknown} thing - The structured data object to serialise.
 * @returns {string} An application/ld+json script tag containing the serialised data.
 */
export const serializeSchema = (thing: unknown): string => {
	return `<script type="application/ld+json">${JSON.stringify(thing)}</script>`;
};
