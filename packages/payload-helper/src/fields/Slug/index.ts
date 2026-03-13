import { type CheckboxField, type TextField, deepMerge } from 'payload';

import { formatSlugHook } from './formatSlug';

type Overrides = {
	slugOverrides?: Partial<TextField>;
	checkboxOverrides?: Partial<CheckboxField>;
	description?: string;
};

type Slug = (fieldToUse?: string, overrides?: Overrides) => [TextField, CheckboxField];

export const SlugField: Slug = (fieldToUse = 'title', overrides = {}) => {
	const { slugOverrides, checkboxOverrides } = overrides;

	const checkBoxField = deepMerge<CheckboxField, Partial<CheckboxField>>(
		{
			name: 'slugLock',
			type: 'checkbox',
			defaultValue: true,
			admin: {
				hidden: true,
				position: 'sidebar',
			},
		},
		checkboxOverrides || {},
	);

	const slugField = deepMerge<TextField, Partial<TextField>>(
		{
			name: 'slug',
			type: 'text',
			index: true,
			label: 'Slug',
			unique: true,
			required: true,
			hooks: {
				beforeValidate: [formatSlugHook(fieldToUse)],
			},
			admin: {
				position: 'sidebar',
				description:
					overrides.description ||
					'The URL friendly version of the title, users will see this text in the URL bar.',
				components: {
					Field: {
						path: '/fields/Slug/Component#Component',
						clientProps: {
							fieldToUse,
							checkboxFieldPath: checkBoxField.name,
						},
					},
				},
			},
		},
		slugOverrides || {},
	);

	return [slugField, checkBoxField];
};

export type { Overrides, Slug };
