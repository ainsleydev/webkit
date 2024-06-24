import type { Field, GlobalConfig, Tab } from 'payload/types';

/**
 * Navigation Configuration for the header and footer
 * nav links.
 */
interface NavigationConfig {
	includeFooter?: boolean;
	header?: NavigationMenuConfig;
	footer?: NavigationMenuConfig;
	additionalTabs?: Tab[];
}

/**
 * Navigation Menu Configuration defines the config for an
 * individual navigation menu, i.e Header or Footer.
 */
interface NavigationMenuConfig {
	maxDepth?: number;
	additionalFields?: Field[];
}

/**
 * Function to generate children structure recursively.
 *
 * @param depth Current depth
 * @param maxDepth Maximum depth to generate children
 * @param fields Fields to include at each level
 */
const generateChildren = (depth: number, maxDepth: number, fields: Field[]): Field[] => {
	if (depth >= maxDepth) {
		return [];
	}

	// Only generate children if depth is less than maxDepth
	if (depth < maxDepth - 1) {
		return [
			{
				name: 'children',
				type: 'array',
				label: `Children Level ${depth + 1}`,
				fields: [...fields, ...generateChildren(depth + 1, maxDepth, fields)],
			},
		];
	}

	return [];
};

/**
 * The default navigation field links.
 */
const navFields: Field[] = [
	{
		type: 'row',
		fields: [
			{
				name: 'title',
				type: 'text',
				label: 'Title',
				required: true,
				admin: {
					width: '50%',
				},
			},
			{
				name: 'url',
				type: 'text',
				label: 'URL',
				required: true,
				admin: {
					width: '50%',
				},
			},
		],
	},
];

/**
 * Navigation Global Configuration
 * Additional fields will be appended to each navigation item.
 *
 * @constructor
 * @param config
 */
export const Navigation = (config?: NavigationConfig): GlobalConfig => {
	const tabs: Tab[] = [
		{
			label: 'Header',
			fields: [
				{
					name: 'header',
					type: 'array',
					label: 'Items',
					interfaceName: 'NavigationHeaderLinks',
					maxRows: 8,
					labels: {
						singular: 'Link',
						plural: 'Links',
					},
					admin: {
						initCollapsed: true,
						isSortable: true,
					},
					fields: [
						...navFields,
						...(config?.header?.maxDepth
							? generateChildren(0, config.header.maxDepth, navFields)
							: []),
						...(config?.header?.additionalFields ? config.header.additionalFields : []),
					],
				},
			],
		},
	];

	if (config?.includeFooter) {
		tabs.push({
			label: 'Footer',
			fields: [
				{
					name: 'footer',
					type: 'array',
					label: 'Items',
					interfaceName: 'NavigationFooterLinks',
					maxRows: 8,
					labels: {
						singular: 'Link',
						plural: 'Links',
					},
					admin: {
						initCollapsed: true,
						isSortable: true,
					},
					fields: [
						...navFields,
						...(config?.footer?.maxDepth
							? generateChildren(0, config.footer.maxDepth, navFields)
							: []),
						...(config?.footer?.additionalFields ? config.footer.additionalFields : []),
					],
				},
			],
		});
	}

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
		fields: [
			{
				type: 'tabs',
				tabs: [...tabs, ...(config.additionalTabs ? config.additionalTabs : [])],
			},
		],
	};
};
