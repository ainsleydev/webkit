import { z } from 'zod';

/**
 * Represents a form field from Payload CMS form builder.
 * This is a simplified type that covers the most common field types.
 */
export type PayloadFormField =
	| {
			name: string;
			label?: string | null;
			placeholder?: string | null;
			defaultValue?: string | null;
			required?: boolean | null;
			id?: string | null;
			blockName?: string | null;
			blockType: 'text';
	  }
	| {
			name: string;
			label?: string | null;
			placeholder?: string | null;
			defaultValue?: string | null;
			required?: boolean | null;
			id?: string | null;
			blockName?: string | null;
			blockType: 'email';
	  }
	| {
			name: string;
			label?: string | null;
			placeholder?: string | null;
			defaultValue?: number | null;
			required?: boolean | null;
			id?: string | null;
			blockName?: string | null;
			blockType: 'number';
	  }
	| {
			name: string;
			label?: string | null;
			placeholder?: string | null;
			defaultValue?: string | null;
			required?: boolean | null;
			id?: string | null;
			blockName?: string | null;
			blockType: 'textarea';
	  }
	| {
			name: string;
			label?: string | null;
			defaultValue?: boolean | null;
			required?: boolean | null;
			id?: string | null;
			blockName?: string | null;
			blockType: 'checkbox';
	  };

/**
 * Generates a dynamic Zod schema from Payload CMS Form Fields.
 *
 * This utility automatically creates validation schemas for Payload form fields,
 * respecting required flags and field types.
 *
 * @param fields - Array of form fields from Payload CMS form builder.
 * @returns Zod object schema for validating the form.
 *
 * @example
 * ```typescript
 * import { generateFormSchema } from '@ainsleydev/sveltekit-helper/utils/forms'
 *
 * const fields = [
 *   { blockType: 'text', name: 'name', label: 'Name', required: true },
 *   { blockType: 'email', name: 'email', label: 'Email', required: true },
 *   { blockType: 'textarea', name: 'message', label: 'Message', required: false }
 * ]
 *
 * const schema = generateFormSchema(fields)
 * // Returns z.object({
 * //   name: z.string().min(1, { message: 'Name is required' }),
 * //   email: z.string().email({ message: 'Invalid email' }),
 * //   message: z.string().optional()
 * // })
 * ```
 */
export const generateFormSchema = (
	fields?: PayloadFormField[] | null,
): z.ZodObject<Record<string, z.ZodTypeAny>> => {
	if (!fields?.length) return z.object({}) as z.ZodObject<Record<string, z.ZodTypeAny>>;

	const shape: Record<string, z.ZodTypeAny> = {};

	for (const field of fields) {
		if (field.blockType === 'text' || field.blockType === 'textarea') {
			shape[field.name] = field.required
				? z.string().min(1, { message: `${field.label || field.name} is required` })
				: z.string().optional();
		} else if (field.blockType === 'email') {
			shape[field.name] = field.required
				? z.string().email({ message: 'Invalid email' })
				: z.string().email({ message: 'Invalid email' }).optional();
		} else if (field.blockType === 'number') {
			shape[field.name] = field.required ? z.number() : z.number().optional();
		} else if (field.blockType === 'checkbox') {
			shape[field.name] = field.required
				? z.boolean().refine((val) => val === true, {
						message: `${field.label || field.name} must be checked`,
					})
				: z.boolean().optional();
		}
	}

	return z.object(shape) as z.ZodObject<Record<string, z.ZodTypeAny>>;
};
