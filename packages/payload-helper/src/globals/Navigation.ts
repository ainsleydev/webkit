import {
	type ArrayField,
	type Field,
	type GlobalConfig,
	type Tab,
	type TabsField,
	deepMerge,
} from 'payload';

/**
 * Navigation Configuration for the header and footer
 * nav links.
 */
interface NavigationConfig {
	overrides?: Partial<Omit<GlobalConfig, 'slug' | 'fields'>>;
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
				name: 'label',
				label: 'Label',
				type: 'text',
				required: true,
				admin: {
					width: '50%',
					description: 'Enter the text that will appear as a label for the link.',
				},
			},
			{
				name: 'url',
				label: 'URL',
				type: 'text',
				required: true,
				admin: {
					width: '50%',
					description: 'Enter a URL where the link will direct too.',
				},
			},
		],
	},
	{
		name: 'newTab',
		label: 'Open in a new tab?',
		type: 'checkbox',
		defaultValue: false,
		admin: {
			description: 'Check this box if you would like the link to open in a new tab.',
		},
	},
];

/**
 * Navigation Global Configuration
 * Additional fields will be appended to each navigation item.
 *
 * @constructor
 * @param args
 */
export const Navigation = (args?: NavigationConfig): GlobalConfig => {
	const tabs: Tab[] = [
		{
			label: 'Header',
			name: 'header',
			fields: [
				{
					name: 'items',
					label: 'Items',
					type: 'array',
					interfaceName: 'NavigationHeaderLinks',
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
			name: 'footer',
			fields: [
				{
					name: 'items',
					label: 'Items',
					type: 'array',
					interfaceName: 'NavigationFooterLinks',
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

	const defaultConfig = {
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
				tabs: [...tabs, ...(args?.additionalTabs ? args.additionalTabs : [])],
			} as TabsField,
		],
	};

	return deepMerge(defaultConfig, args?.overrides || {});
};
