import path from 'path'
import {slateEditor} from '@payloadcms/richtext-slate'
import type {CollectionConfig} from 'payload/types'

export const Media: CollectionConfig = {
	slug: 'media',
	upload: {
		staticURL: '/media',
		staticDir: path.resolve(__dirname, '../../../media'),
		externalFileHeaderFilter: (headers) => {
			return {
				...headers,
				'CacheControl': 'public, max-age=31536000',
			}
		},
		imageSizes: [
			// Original Size (for WebP & Avif)
			{
				name: 'webp',
				width: undefined,
				height: undefined,
				formatOptions: {
					format: 'webp',
					options: {
						quality: 80,
					},
				}
			},
			{
				name: 'avif',
				width: undefined,
				height: undefined,
				formatOptions: {
					format: 'avif',
					options: {
						quality: 80,
					},
				}
			},
			// Thumbnail Sizes
			{
				name: 'thumbnail',
				width: 400,
				height: 300,
				position: 'centre',
			},
			{
				name: 'thumbnail_webp',
				width: 400,
				height: 300,
				position: 'centre',
				formatOptions: {
					format: 'webp',
					options: {
						quality: 80,
					},
				}
			},
			{
				name: 'thumbnail_avif',
				width: 400,
				height: 300,
				position: 'centre',
				formatOptions: {
					format: 'avif',
					options: {
						quality: 80,
					},
				}
			},
			// Mobile Sizes
			{
				name: 'mobile',
				width: 768,
				height: undefined,
			},
			{
				name: 'mobile_webp',
				width: 768,
				height: undefined,
				formatOptions: {
					format: 'webp',
					options: {
						quality: 80,
					},
				}
			},
			{
				name: 'mobile_avif',
				width: 768,
				height: undefined,
				formatOptions: {
					format: 'avif',
					options: {
						quality: 80,
					},
				}
			},
			// Tablet Sizes
			{
				name: 'tablet',
				width: 1024,
				height: undefined,
			},
			{
				name: 'tablet_webp',
				width: 1024,
				height: undefined,
				formatOptions: {
					format: 'webp',
					options: {
						quality: 80,
					},
				}
			},
			{
				name: 'tablet_avif',
				width: 1024,
				height: undefined,
				formatOptions: {
					format: 'avif',
					options: {
						quality: 80,
					},
				}
			},
		],
	},
	fields: [
		{
			name: 'alt',
			type: 'text',
			required: true,
		},
		{
			name: 'caption',
			type: 'richText',
			editor: slateEditor({
				admin: {
					elements: ['link'],
				},
			}),
		},
	],
}
