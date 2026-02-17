import type { Field, FieldHookArgs, TextField, TypeWithID } from 'payload';
import { deepMerge } from 'payload';

export type URLFieldArgs<T extends TypeWithID> = {
	overrides?: Partial<Omit<TextField, 'type'>>;
	generate: (args: FieldHookArgs<T>) => string | Promise<string | undefined>;
};

/**
 * Creates a virtual URL field with a custom generation function.
 *
 * @param generate - A function that generates the URL based on the field data.
 * @param overrides - Optional overrides to customise the field.
 */
export const URLField = <T extends TypeWithID>({ generate, overrides }: URLFieldArgs<T>): Field => {
	const baseField: Field = {
		name: 'url',
		label: 'URL',
		type: 'text',
		admin: {
			readOnly: true,
			position: 'sidebar',
		},
		virtual: true,
		hooks: {
			afterRead: [
				async (args: FieldHookArgs<T>) => {
					let url = await generate(args);
					if (args.draft) {
						url += '?draft=true';
					}
					return url;
				},
			],
		},
	};
	return deepMerge<Field, Partial<Omit<TextField, 'type'>>>(baseField, overrides || {});
};
