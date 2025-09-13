import { type ArrayField, type Field, type GlobalConfig, type Tab, deepMerge } from 'payload';

/**
 * Navigation arguments for the header and footer
 * nav links.
 */
interface NavigationArgs {
	includeFooter?: boolean;
	header?: NavigationMenuConfig;
	footer?: NavigationMenuConfig;
	additionalTabs?: Tab[];
	overrides?: Partial<GlobalConfig>;
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
 */
export const Navigation = (args?: NavigationArgs): GlobalConfig => {
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
						...(args?.header?.maxDepth
							? generateChildren(0, args.header.maxDepth, navFields)
							: []),
						...(args?.header?.additionalFields ? args.header.additionalFields : []),
					],
				},
			],
		} as ArrayField,
	];

	if (args?.includeFooter) {
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
						...(args?.footer?.maxDepth
							? generateChildren(0, args.footer.maxDepth, navFields)
							: []),
						...(args?.footer?.additionalFields ? args.footer.additionalFields : []),
					],
				},
			],
		} as ArrayField);
	}

	const defaultConfig: GlobalConfig = {
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
				tabs: [...tabs, ...(args?.additionalTabs ?? [])] as Tab[],
			},
		],
	};

	return deepMerge(defaultConfig, args?.overrides || {});
};
