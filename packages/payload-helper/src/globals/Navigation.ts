import type { Field, GlobalConfig } from 'payload/types';

/**
 * Navigation Global Configuration
 * Additional fields will be appended to each navigation item.
 *
 * @param additionalFields
 * @constructor
 */
export const Navigation = (additionalFields?: Field[]): GlobalConfig => {
	return {
		slug: 'navigation',
		typescript: {
			interface: 'Navigation',
		},
		graphQL: {
			name: 'Navigation',
		},
		access: {
			read: () => true,
		},
		fields: [...additionalFields],
	};
};
