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
