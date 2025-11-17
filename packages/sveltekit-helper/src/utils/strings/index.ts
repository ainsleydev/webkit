/**
 * Generates an alphanumeric random string.
 *
 * @param length - The length of the string to generate.
 * @returns The generated random string.
 *
 * @example
 * ```typescript
 * const id = generateRandomString(10)
 * // Returns something like: "a3k9mxp2q1"
 * ```
 */
export const generateRandomString = (length: number): string => {
	let result = '';
	while (result.length < length) {
		result += (Math.random() + 1).toString(36).substring(2);
	}
	return result.substring(0, length);
};
