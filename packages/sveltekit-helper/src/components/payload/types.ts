import type { HTMLAttributes } from 'svelte/elements';

export type MediaSizes = Record<
	string,
	| Partial<{
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
	  }>
	| undefined
>;

export type PayloadMediaProps = HTMLAttributes<HTMLElement> & {
	data: Media;
	loading?: 'lazy' | 'eager' | undefined;
	className?: string;
	breakpointBuffer?: number;
	maxWidth?: number | undefined;
	onload?: (event: Event) => void;
};

export type Media = {
	id: number;
	alt?: string;
	updatedAt: string;
	createdAt: string;
	deletedAt?: string | null;
	url?: string | null;
	thumbnailURL?: string | null;
	filename?: string | null;
	mimeType?: string | null;
	filesize?: number | null;
	width?: number | null;
	height?: number | null;
	focalX?: number | null;
	focalY?: number | null;
	sizes?: MediaSizes;
};

/**
 * Props that the PayloadSEO component expects. Accepts site-level settings
 * and an optional page meta object which is merged with higher priority.
 */
export type PayloadSEOProps = {
	siteName: string;
	settings: PayloadSettings;
	pageMeta?: PayloadMeta;
	pageCodeInjection?: PayloadCodeInjection;
};

/**
 * Props that the PayloadFooter component expects.
 * Renders code injection scripts for both settings and page level.
 */
export type PayloadFooterProps = {
	settings: PayloadSettings;
	pageCodeInjection?: PayloadCodeInjection;
};

/**
 * Meta exported from the Payload SEO plugin. Appears on
 * both page and settings level.
 */
export type PayloadMeta = {
	title?: string | null;
	description?: string | null;
	image?: (number | null) | Media;
	private?: boolean | null;
	canonicalURL?: string | null;
	structuredData?: unknown;
};

/**
 * Settings compatible with the type generated in Payload
 * under the webkit settings global.
 */
export type PayloadSettings = {
	siteName?: string | null;
	locale?: string;
	tagLine?: string | null;
	contact?: {
		email?: string | null;
		telephone?: string | null;
	};
	address?: {
		line1?: string | null;
		line2?: string | null;
		city?: string | null;
		county?: string | null;
		postcode?: string | null;
		country?: string | null;
	};
	meta?: PayloadMeta;
	social?: PayloadSocial;
	codeInjection?: PayloadCodeInjection;
};

/**
 * Social links that appear in the organisation settings.
 */
export type PayloadSocial = {
	linkedIn?: string | null;
	x?: string | null;
	facebook?: string | null;
	instagram?: string | null;
	youtube?: string | null;
	tiktok?: string | null;
};

/**
 * Header and footer scripts to be injected.
 */
export type PayloadCodeInjection = {
	head?: string | null;
	footer?: string | null;
};
