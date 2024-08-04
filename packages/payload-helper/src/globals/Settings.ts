import type { GlobalConfig, GroupField, Tab, UploadField } from 'payload';
import { validatePostcode, validateURL } from '../util/validation.js';
import { countries } from './countries.js';
import { languages } from './locales.js';

/**
 * Settings Global Configuration
 * Additional tabs will be appended to the settings page.
 * TODO, type error in here somewhere.
 *
 * @param additionalTabs
 * @constructor
 */
export const Settings = (additionalTabs: Tab[]): GlobalConfig => {
	return {
		slug: 'settings',
		typescript: {
			interface: 'Settings',
		},
		graphQL: {
			name: 'Settings',
		},
		access: {
			read: () => true,
		},
		fields: [
			{
				type: 'tabs',
				tabs: [
					{
						label: 'Global',
						description: 'Configure global settings for the website.',
						fields: [
							{
								type: 'row',
								fields: [
									{
										name: 'siteName',
										type: 'text',
										label: 'Site Name',
										admin: {
											width: '50%',
											description:
												'Add a site name for the website, this will be outputted in the Open Graph schema as well as a suffix for the meta title.',
										},
									},
									{
										name: 'locale',
										type: 'select',
										label: 'Locale',
										defaultValue: 'en_GB',
										options: languages.map((l) => {
											return {
												label: l.name,
												value: l.code,
											};
										}),
										admin: {
											width: '50%',
											description:
												'Add a locale for the website, this will be outputted in the Open Graph schema and the top level HTML tag. Defaults to en_GB.',
										},
										typescriptSchema: [
											() => ({
												type: 'string',
											}),
										],
									},
								],
							},
							{
								name: 'tagLine',
								type: 'textarea',
								label: 'Tag Line',
								admin: {
									description: 'In a few words, explain what this site is about',
								},
							},
							{
								name: 'logo',
								type: 'upload',
								relationTo: 'media',
								filterOptions: {
									mimeType: {
										contains: 'image',
									},
								},
								admin: {
									description:
										'Add a logo for the website that will be displayed in the header & across the website.',
								},
							} as UploadField,
							{
								name: 'robots',
								type: 'textarea',
								label: 'Robots.txt',
								admin: {
									description:
										'Robots.txt is a text file webmasters create to instruct web robots (typically search engine robots) how to crawl pages on their website.',
								},
							},
						],
					},
					{
						label: 'Code Injection',
						description:
							'Code injection allows you to inject a small snippet of HTML into your site. It can be a css override, analytics of a block javascript.',
						fields: [
							{
								name: 'codeInjection',
								type: 'group',
								fields: [
									{
										name: 'head',
										type: 'code',
										label: 'Head',
										admin: {
											language: 'html',
											description:
												'Outputs code within the <head> of the website.',
										},
									},
									{
										name: 'footer',
										type: 'code',
										label: 'Footer',
										admin: {
											language: 'html',
											description:
												'Outputs code in the footer of the website.',
										},
									},
								],
							},
						],
					},
					{
						label: 'Contact Details',
						fields: [
							{
								name: 'contact',
								type: 'group',
								admin: {
									hideGutter: true,
									description:
										'Add global contact details for the website that will be used in schema & contact pages.',
								},
								fields: [
									{
										type: 'row',
										fields: [
											{
												name: 'email',
												type: 'email',
												label: 'Email',
												admin: {
													width: '50%',
												},
											},
											{
												name: 'telephone',
												type: 'text',
												label: 'Telephone',
												admin: {
													width: '50%',
												},
											},
										],
									},
								],
							} as GroupField,
							{
								type: 'group',
								name: 'address',
								label: 'Address',
								admin: {
									hideGutter: true,
									description: 'Add an address for the website.',
								},
								fields: [
									{
										type: 'row',
										fields: [
											{
												name: 'line1',
												type: 'text',
												label: 'Line 1',
												admin: {
													width: '50%',
												},
											},
											{
												name: 'line2',
												type: 'text',
												label: 'Line 2',
												admin: {
													width: '50%',
												},
											},
											{
												name: 'city',
												type: 'text',
												label: 'City',
												admin: {
													width: '50%',
												},
											},
											{
												name: 'county',
												type: 'text',
												label: 'County',
												admin: {
													width: '50%',
												},
											},
											{
												name: 'postcode',
												type: 'text',
												label: 'Postcode',
												validate: validatePostcode,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'country',
												type: 'select',
												label: 'Country',
												options: countries.map((c) => {
													return {
														label: c,
														value: c,
													};
												}),
												admin: {
													width: '50%',
												},
											},
										],
									},
								],
							} as GroupField,
							{
								type: 'group',
								name: 'social',
								label: 'Social Links',
								admin: {
									hideGutter: true,
									description: 'Add social links for the website.',
								},
								fields: [
									{
										type: 'row',
										fields: [
											{
												name: 'linkedIn',
												type: 'text',
												label: 'LinkedIn',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'x',
												type: 'text',
												label: 'X',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'facebook',
												type: 'text',
												label: 'Facebook',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'instagram',
												type: 'text',
												label: 'Instagram',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'youtube',
												type: 'text',
												label: 'Youtube',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
											{
												name: 'tiktok',
												type: 'text',
												label: 'TikTok',
												validate: validateURL,
												admin: {
													width: '50%',
												},
											},
										],
									},
								],
							} as GroupField,
						],
					},
					{
						label: 'Maintenance',
						fields: [
							{
								name: 'maintenance',
								type: 'group',
								fields: [
									{
										name: 'enabled',
										type: 'checkbox',
										label: 'Enable',
										admin: {
											description:
												'Enable maintenance mode for the site, this will use a maintenance page template and not include any of the sites functioanlity.',
										},
									},
									{
										name: 'title',
										type: 'text',
										label: 'Title',
										admin: {
											description: 'Add a title for the maintenance page.',
										},
									},
									{
										name: 'content',
										type: 'textarea',
										label: 'Content',
										admin: {
											description:
												'Add content for the maintenance page, it will appear beneath the title.',
										},
									},
								],
							},
						],
					},
					...(additionalTabs ? additionalTabs : []),
				],
			},
		],
	};
};
