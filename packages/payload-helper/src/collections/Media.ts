import type { CollectionConfig, Field } from 'payload';
// import type {
// 	LexicalRichTextAdapterProvider,
// 	LexicalEditorProps,
// } from '@payloadcms/richtext-lexical';

/**
 * Media Collection Configuration
 * Additional fields will be appended to the media collection.
 *
 * @constructor
 * @param editor
 * @param additionalFields
 */
export const Media = (
	//editor?: (props?: LexicalEditorProps) => LexicalRichTextAdapterProvider,
	additionalFields?: Field[],
): CollectionConfig => {
	return {
		slug: 'media',
		access: {
			read: () => true,
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
				required: false,
				// editor: editor({
				// 	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
				// 	// @ts-ignore
				// 	features: ({ defaultFeatures }) => {
				// 		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
				// 		// @ts-ignore
				// 		return defaultFeatures.filter((feature) => {
				// 			return feature.key === 'paragraph' || feature.key === 'link';
				// 		});
				// 	},
				// }),
			},
			...(additionalFields ? additionalFields : []),
		],
		upload: {
			staticDir: 'media',
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
				// Desktop Sizes
				{
					name: 'desktop',
					width: 1440,
					height: undefined,
				},
				{
					name: 'desktop_webp',
					width: 1440,
					height: undefined,
					formatOptions: {
						format: 'webp',
						options: {
							quality: 80,
						},
					},
				},
				{
					name: 'desktop_avif',
					width: 1440,
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
	};
};
