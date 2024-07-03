import type { CollectionConfig } from 'payload';

/**
 * Redirects Collection Configuration
 * In favour of the native plugin for more granular control.
 *
 * TODO: Potential to add regex redirects here, i.e product-category/(.*)$ => /category/$1
 *
 * @constructor
 */
export const Redirects = (overrides?: Partial<CollectionConfig>): CollectionConfig => {
	return {
		slug: 'redirects',
		admin: {
			useAsTitle: 'from',
		},
		fields: [
			{
				name: 'from',
				type: 'text',
				label: 'From URL',
				required: true,
				index: true,
				admin: {
					description: 'The URL you want to redirect from, ensure it starts with a /',
				},
			},
			{
				name: 'to',
				type: 'text',
				label: 'Destination URL',
				required: true,
				admin: {
					description:
						'The URL you want to redirect to, can be a relative or absolute URL',
				},
			},
			{
				name: 'code',
				type: 'select',
				label: 'Redirect Code',
				required: true,
				defaultValue: '301',
				options: [
					{ label: '301 - Permanent', value: '301' },
					{ label: '302 - Temporary', value: '302' },
					{ label: '307 - Temporary Redirect', value: '307' },
					{ label: '308 - Permanent Redirect', value: '308' },
					{ label: '410 - Content Deleted', value: '410' },
					{ label: '451 - Unavailable For Legal Reasons', value: '451' },
				],
			},
		],
		...(overrides ? overrides : {}),
	};
};
