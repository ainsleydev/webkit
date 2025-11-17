import type { z } from 'zod';

export type FormErrors = Record<string, string>;

/**
 * Flattens Zod validation errors into a simple key-value object.
 * Takes only the first error message for each field.
 *
 * @param error - ZodError from validation failure.
 * @returns Object mapping field names to error messages.
 *
 * @example
 * ```typescript
 * import { flattenZodErrors } from '@ainsleydev/sveltekit-helper/utils/forms'
 *
 * const schema = z.object({ email: z.string().email() })
 * const result = schema.safeParse({ email: 'invalid' })
 *
 * if (!result.success) {
 *   const errors = flattenZodErrors(result.error)
 *   // { email: 'Invalid email' }
 * }
 * ```
 */
export const flattenZodErrors = (error: z.ZodError): FormErrors => {
	const errors: FormErrors = {};

	for (const issue of error.issues) {
		const path = issue.path.join('.');
		if (!errors[path]) {
			errors[path] = issue.message;
		}
	}

	return errors;
};
