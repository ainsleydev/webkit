import { enhance as kitEnhance } from '$app/forms';
import { type Readable, type Writable, get, writable } from 'svelte/store';
import { z } from 'zod';

import type { SubmitFunction } from '@sveltejs/kit';

export type FormData = Record<string, never>;
export type FormErrors = Record<string, string>;

export type ClientForm = {
	fields: Writable<FormData>;
	errors: Writable<FormErrors>;
	validate: (opts?: { field?: string }) => boolean;
	submitting: Readable<boolean>;
	submit: (submitter?: HTMLElement | Event | EventTarget | null) => void;
	submitted: Readable<boolean>;
	enhance: (el: HTMLFormElement, submitFn?: SubmitFunction) => { destroy(): void };
};

export type ServerForm<T = FormData> = {
	valid: boolean;
	data: T;
	errors: FormErrors;
};

/**
 * Creates a reactive client-side form store with validation and submission handling.
 *
 * This utility combines Svelte stores with Zod validation to provide a comprehensive
 * form management solution for SvelteKit applications.
 *
 * @param schema - Zod schema to validate form fields against.
 * @param options - Client form options.
 * @param options.submissionDelay - Optional delay (in milliseconds) before setting submitting to false.
 * @param onSubmit - Optional async submit handler called after successful validation.
 * @param initial - Optional initial form values.
 * @returns ClientForm object with reactive stores and methods.
 *
 * @example
 * ```typescript
 * import { z } from 'zod'
 * import { clientForm } from '@ainsleydev/sveltekit-helper/utils/forms'
 *
 * const schema = z.object({
 *   email: z.string().email(),
 *   password: z.string().min(8)
 * })
 *
 * const { fields, errors, validate, submitting, enhance } = clientForm(
 *   schema,
 *   { submissionDelay: 300 },
 *   async (data) => {
 *     await fetch('/api/login', {
 *       method: 'POST',
 *       body: JSON.stringify(data)
 *     })
 *   }
 * )
 * ```
 */
export const clientForm = (
	schema: z.ZodSchema,
	options?: { submissionDelay?: number },
	onSubmit?: (data: FormData) => Promise<void>,
	initial: FormData = {},
): ClientForm => {
	const fields = writable<FormData>(initial);
	const errors = writable<FormErrors>({});
	const submitting = writable<boolean>(false);
	const submitted = writable<boolean>(false);
	const submissionDelay = options?.submissionDelay ?? 0;

	const setSubmittingFalse = () => {
		if (submissionDelay > 0) {
			setTimeout(() => submitting.set(false), submissionDelay);
		} else {
			submitting.set(false);
		}
	};

	/**
	 * Validates the entire form or a specific field.
	 *
	 * @param opts - Validation options.
	 * @param opts.field - Optional field name to validate only that field.
	 * @returns True if validation passed, false otherwise.
	 */
	const validate = (opts?: { field?: string }) => {
		const current = get(fields);

		if (opts?.field && schema instanceof z.ZodObject) {
			const fieldName: string = opts.field;
			const fieldSchema = schema.shape[fieldName];
			if (!fieldSchema) return true;

			const result = fieldSchema.safeParse(current[fieldName]);
			errors.update((prev) => {
				const rest = { ...prev };
				if (!result.success) {
					rest[fieldName] = result.error.issues[0]?.message || 'Invalid value';
				} else {
					delete rest[fieldName];
				}
				return rest;
			});

			return result.success;
		}

		const result = schema.safeParse(current);
		errors.set(result.success ? {} : flattenZodErrors(result.error));
		return result.success;
	};

	/**
	 * Submits the form programmatically.
	 * Typically used with on:submit={submit} handler.
	 *
	 * @param submitter - Optional HTMLElement, Event, or EventTarget that triggered submission.
	 */
	const submit = async (submitter?: HTMLElement | Event | EventTarget | null) => {
		if (submitter instanceof Event) {
			submitter.preventDefault();
		}

		submitted.set(true);

		if (!validate()) return;
		if (!onSubmit) return;

		submitting.set(true);
		try {
			await onSubmit(get(fields));
		} finally {
			setSubmittingFalse();
		}
	};

	/**
	 * Enhances an HTML form element with SvelteKit form actions.
	 * Wraps SvelteKit's enhance with validation and custom submission logic.
	 *
	 * @param el - The HTML form element to enhance.
	 * @param submitFn - Optional custom submit function.
	 * @returns Object with destroy method to clean up the enhancement.
	 */
	const enhance = (el: HTMLFormElement, submitFn?: SubmitFunction) => {
		const defaultSubmitFn: SubmitFunction = ({ cancel }) => {
			submitted.set(true);
			submitting.set(true);

			if (!validate()) {
				setSubmittingFalse();
				cancel();
				return;
			}

			if (onSubmit) {
				cancel();
				onSubmit(get(fields)).finally(() => {
					setSubmittingFalse();
				});
				return;
			}

			return async ({ result, update }) => {
				submitting.set(false);

				if (result.type === 'success') {
					// Handle success
				}

				await update();
			};
		};

		return kitEnhance(el, submitFn || defaultSubmitFn);
	};

	return {
		fields,
		errors,
		validate,
		submit,
		submitting,
		submitted,
		enhance,
	};
};

/**
 * Validates form data on the server using a Zod schema.
 * Parses both JSON and FormData submissions.
 *
 * @param input - Request object or pre-parsed data object.
 * @param schema - Zod schema used to validate the input data.
 * @returns Promise resolving to ServerForm with validation results.
 *
 * @example
 * ```typescript
 * import { serverForm } from '@ainsleydev/sveltekit-helper/utils/forms'
 * import { z } from 'zod'
 *
 * export const actions = {
 *   default: async ({ request }) => {
 *     const schema = z.object({
 *       email: z.string().email(),
 *       password: z.string().min(8)
 *     })
 *
 *     const { valid, data, errors } = await serverForm(request, schema)
 *
 *     if (!valid) {
 *       return { errors }
 *     }
 *
 *     // Process form...
 *   }
 * }
 * ```
 */
export const serverForm = async <T = FormData>(
	input: Request | unknown,
	schema: z.ZodSchema,
): Promise<ServerForm<T>> => {
	let data: unknown;

	if (input instanceof Request) {
		const contentType = input.headers.get('content-type') || '';

		if (contentType.includes('application/json')) {
			data = await input.json();
		} else {
			const formData = await input.formData();
			data = Object.fromEntries(formData.entries());
		}
	} else {
		data = input;
	}

	if (!schema) {
		return {
			valid: true,
			data: data as T,
			errors: {},
		};
	}

	const result = schema.safeParse(data);
	return {
		valid: result.success,
		data: data as T,
		errors: result.success ? {} : flattenZodErrors(result.error),
	};
};

/**
 * Converts a ZodError into a simple key-to-message error map.
 * Takes only the first error message for each field.
 *
 * @param error - ZodError from validation failure.
 * @returns Object mapping field names to error messages.
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
