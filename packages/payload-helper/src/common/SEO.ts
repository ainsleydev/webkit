import { Field } from 'payload/types';
import { validateURL } from '../util/validation';

/**
 * SEO Fields define the additional fields that appear
 * within the Payload SEO plugin.
 */
export const SEOFields: Field[] = [
	{
		type: 'row',
		fields: [
			{
				name: 'private',
				type: 'checkbox',
				label: 'Private',
				defaultValue: false,
				admin: {
					width: '50%',
					description:
						'Enable private mode to prevent robots from crawling the page or website. When enabled it will output <meta name="robots" content="noindex" /> on the frontend.',
				},
			},
			{
				name: 'canonicalURL',
				type: 'text',
				label: 'Canonical',
				admin: {
					width: '50%',
					description:
						'A canonical URL is the version of a webpage chosen by search engines like Google as the main version when there are duplicates.',
				},
				validate: validateURL,
			},
		],
	},
	{
		name: 'structuredData',
		type: 'json',
		label: 'Structured Data',
		admin: {
			description:
				'Structured data is a standardized format for providing information about a page and classifying the page content. The site Schema.org contains a standardised list of markup that the major search engines — Google, Bing, Yahoo and Yandex — have collectively agreed to support.',
		},
	},
];
