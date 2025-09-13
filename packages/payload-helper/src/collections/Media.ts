import { lexicalEditor } from '@payloadcms/richtext-lexical';
import * as mime from 'mime-types';
import type { CollectionConfig, Field, PayloadRequest, UploadConfig } from 'payload';
import type { ImageSize } from 'payload';
import env from '../util/env.js';

export interface MediaArgs {
	includeAvif?: boolean;
	additionalFields?: Field[];
	uploadOverrides?: Partial<UploadConfig>;
}

export const imageSizes: ImageSize[] = [
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
		name: 'thumbnail',
		width: 200,
		height: undefined,
		position: 'centre',
	},
	{
		name: 'thumbnail_webp',
		width: 200,
		height: undefined,
		position: 'centre',
		formatOptions: {
			format: 'webp',
			options: {
				quality: 80,
			},
		},
	},
	{
		name: 'mobile',
		width: 500,
		height: undefined,
	},
	{
		name: 'mobile_webp',
		width: 500,
		height: undefined,
		formatOptions: {
			format: 'webp',
			options: {
				quality: 80,
			},
		},
	},
	{
		name: 'tablet',
		width: 800,
		height: undefined,
	},
	{
		name: 'tablet_webp',
		width: 800,
		height: undefined,
		formatOptions: {
			format: 'webp',
			options: {
				quality: 80,
			},
		},
	},
	{
		name: 'desktop',
		width: 1200,
		height: undefined,
	},
	{
		name: 'desktop_webp',
		width: 1200,
		height: undefined,
		formatOptions: {
			format: 'webp',
			options: {
				quality: 80,
			},
		},
	},
];

export const imageSizesWithAvif = (): ImageSize[] => {
	return [
		...imageSizes,
		{
			name: 'avif',
			width: undefined,
			height: undefined,
			formatOptions: {
				format: 'avif',
				options: {
					quality: 60,
					effort: 1,
					chromaSubsampling: '4:4:4',
					bitdepth: 8,
					lossless: false,
				},
			},
		},
		{
			name: 'thumbnail_avif',
			width: 200,
			height: undefined,
			position: 'centre',
			formatOptions: {
				format: 'avif',
				options: {
					quality: 60,
					effort: 1,
					chromaSubsampling: '4:4:4',
					bitdepth: 8,
					lossless: false,
				},
			},
		},
		{
			name: 'mobile_avif',
			width: 500,
			height: undefined,
			formatOptions: {
				format: 'avif',
				options: {
					quality: 60,
					effort: 1,
					chromaSubsampling: '4:4:4',
					bitdepth: 8,
					lossless: false,
				},
			},
		},
		{
			name: 'tablet_avif',
			width: 800,
			height: undefined,
			formatOptions: {
				format: 'avif',
				options: {
					quality: 60,
					effort: 1,
					chromaSubsampling: '4:4:4',
					bitdepth: 8,
					lossless: false,
				},
			},
		},
		{
			name: 'desktop_avif',
			width: 1200,
			height: undefined,
			formatOptions: {
				format: 'avif',
				options: {
					quality: 80,
				},
			},
		},
	];
};

export const Media = (args: MediaArgs = {}): CollectionConfig => {
	let sizes = imageSizes;
	if (args.includeAvif) {
		sizes = imageSizesWithAvif();
	}

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
				editor: lexicalEditor({
					features: ({ defaultFeatures }) => {
						return defaultFeatures.filter((feature) => {
							return feature.key === 'paragraph' || feature.key === 'link';
						});
					},
				}),
			},
			...(args.additionalFields ? args.additionalFields : []),
		],
		upload: {
			staticDir: 'media',
			adminThumbnail: 'tablet',
			disableLocalStorage: env.isProduction,
			handlers: [
				async (req: PayloadRequest, args) => {
					const logger = req.payload.logger;
					const { params } = args;
					const { collection, filename } = params;

					if (collection !== 'media') {
						return;
					}

					const contentType = mime.lookup(filename);
					if (!contentType) {
						logger.error(`Unable to find mime type for file: ${filename}`);
						return;
					}

					const headers = new Headers();
					headers.set('Content-Type', contentType);
					headers.set('Cache-Control', 'public, max-age=31536000');

					req.responseHeaders = headers;
				},
			],
			imageSizes: sizes,
			...(args.uploadOverrides ? args.uploadOverrides : {}),
		},
	};
};
