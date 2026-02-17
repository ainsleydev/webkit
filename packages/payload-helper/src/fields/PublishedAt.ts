import type { DateField, Field } from 'payload';

export type PublishedAtArgs = {
	overrides?: Partial<DateField>;
};

/**
 * Creates a published at date field with sensible defaults.
 *
 * Automatically sets the current date as default value and populates
 * the field when a document is first published.
 *
 * @param args - Optional arguments to customise the field.
 */
export const PublishedAt = (args?: PublishedAtArgs): Field => {
	return {
		name: 'publishedAt',
		type: 'date',
		required: true,
		defaultValue: () => new Date().toISOString(),
		admin: {
			position: 'sidebar',
			date: {
				pickerAppearance: 'dayOnly',
			},
			...args?.overrides?.admin,
		},
		hooks: {
			beforeChange: [
				({ siblingData, value }) => {
					if (siblingData._status === 'published' && !value) {
						return new Date();
					}
					return value;
				},
			],
		},
		...args?.overrides,
	};
};
