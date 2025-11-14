export interface MediaSeed {
	path: string;
	alt: string;
	caption?: string;
}

export interface Media {
	id: number;
	alt: string;
	caption?: {
		root: {
			type: string;
			children: {
				type: string;
				version: number;
				[k: string]: unknown;
			}[];
			direction: ('ltr' | 'rtl') | null;
			format: 'left' | 'start' | 'center' | 'right' | 'end' | 'justify' | '';
			indent: number;
			version: number;
		};
		[k: string]: unknown;
	} | null;
	updatedAt: string;
	createdAt: string;
	url?: string | null;
	thumbnailURL?: string | null;
	filename?: string | null;
	mimeType?: string | null;
	filesize?: number | null;
	width?: number | null;
	height?: number | null;
	focalX?: number | null;
	focalY?: number | null;
	sizes?: {
		webp?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		avif?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		thumbnail?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		thumbnail_webp?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		thumbnail_avif?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		mobile?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		mobile_webp?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		mobile_avif?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		tablet?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		tablet_webp?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
		tablet_avif?: {
			url?: string | null;
			width?: number | null;
			height?: number | null;
			mimeType?: string | null;
			filesize?: number | null;
			filename?: string | null;
		};
	};
}
