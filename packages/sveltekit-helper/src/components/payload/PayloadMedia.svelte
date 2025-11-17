<script lang="ts">
import type { Media, MediaSizes, PayloadMediaProps } from './types.js';

const {
	data,
	loading = undefined,
	className = '',
	breakpointBuffer = 50,
	maxWidth = undefined,
	onload = () => null,
	...restProps
}: PayloadMediaProps = $props();

/**
 * Returns sources grouped by breakpoint/width, then sorted by mime priority within each group:
 * For each width: avif → webp → jpeg/png → others
 */
const getSources = (sizesMap?: MediaSizes) => {
	if (!sizesMap) return [];

	const arr = Object.values(sizesMap)
		.filter((s) => s?.url && s?.width)
		.map((s) => ({
			url: s.url as string,
			width: s.width as number,
			mimeType: s.mimeType ?? '',
		}));

	const filtered = maxWidth ? arr.filter((s) => s.width <= maxWidth) : arr;

	const mimePriority = (mime: string, url: string) => {
		if (/avif/i.test(mime) || /\.avif$/i.test(url)) return 1;
		if (/webp/i.test(mime) || /\.webp$/i.test(url)) return 2;
		if (/jpe?g|png/i.test(mime) || /\.(jpe?g|png)$/i.test(url)) return 3;
		return 4;
	};

	return filtered.sort((a, b) => {
		if (a.width !== b.width) return a.width - b.width;
		return mimePriority(a.mimeType, a.url) - mimePriority(b.mimeType, b.url);
	});
};

const isImage = $derived((data?.mimeType ?? '').startsWith('image'));
const isVideo = $derived((data?.mimeType ?? '').startsWith('video'));
const isSVG = $derived(!!(data?.filename && /\.svg$/i.test(data.filename)));

const sources = $derived(isImage && !isSVG ? getSources(data.sizes) : []);
const fallbackSource = $derived(sources[sources.length - 1]);
const imgAlt = $derived(data?.alt ?? '');
const imgLoading = $derived(loading ?? undefined);
const imgWidth = $derived(fallbackSource?.width ?? data?.width ?? undefined);
const imgHeight = $derived(
	fallbackSource && data?.width && data?.height
		? Math.round((fallbackSource.width / data.width) * data.height)
		: data?.height ?? undefined,
);
const fallbackUrl = $derived(sources[sources.length - 1]?.url ?? data?.url ?? '');
</script>

<!--
	@component

	PayloadMedia renders responsive images and videos from Payload CMS media fields.
	Handles multiple image sizes with automatic AVIF/WebP format prioritisation.

	@example
	```svelte
	<PayloadMedia
		data={media}
		loading="lazy"
		maxWidth={1200}
	/>
	```

	@example
	```svelte
	<PayloadMedia
		data={videoMedia}
		className="custom-video"
	/>
	```
-->
{#if isImage}
	{#if isSVG}
		<img
			src={data.url}
			alt={imgAlt}
			width={imgWidth}
			height={imgHeight}
			loading={imgLoading}
			class={className}
			{onload}
			{...restProps}
		/>
	{:else}
		<picture class={className} {...restProps}>
			{#each sources as source, index (index)}
				<source
					media={`(max-width: ${source.width + breakpointBuffer}px)`}
					srcset={source.url}
					type={source.mimeType}
				/>
			{/each}
			<img
				src={fallbackUrl}
				alt={imgAlt}
				width={imgWidth}
				height={imgHeight}
				loading={imgLoading}
				{onload}
			/>
		</picture>
	{/if}
{:else if isVideo}
	<video
		controls
		width={imgWidth}
		height={imgHeight}
		preload={loading === 'lazy' ? 'metadata' : 'auto'}
		poster={data.thumbnailURL}
		class={className}
		{onload}
		{...restProps}
	>
		<source src={data.url} type={data.mimeType} />
		<track kind="captions" />
	</video>
{/if}

<style lang="scss">
	img,
	video {
		max-width: 100%;
		height: auto;
		display: block;
	}
</style>
