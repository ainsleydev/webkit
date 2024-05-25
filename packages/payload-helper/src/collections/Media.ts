import path from 'path';
import { slateEditor } from '@payloadcms/richtext-slate';
import type { CollectionConfig, Field } from 'payload/types';

/**
 * Media Collection Configuration
 * Additional fields will be appended to the media collection.
 *
 * @param additionalFields
 * @constructor
 */
export const Media = (additionalFields?: Field[]): CollectionConfig => {
	return {
		slug: 'media',
		upload: {
			staticURL: '/media',
			staticDir: path.resolve(__dirname, '../../../media'),
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
					},
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
					},
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
					},
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
					},
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
					},
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
					},
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
					},
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
					},
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
			...(additionalFields ? additionalFields : []),
		],
	};
};
