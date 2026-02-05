<script lang="ts">
import { serializeSchema } from '../../utils/seo/ld-json.js';
import { resolveItems } from './resolve.js';

import type { PayloadSEOProps } from './types.js';

const { siteName, pageMeta, settings }: PayloadSEOProps = $props();

/**
 * Meta properties that should appear on both the
 * page and global level.
 */
const meta = $derived.by(() => {
	return {
		title: pageMeta?.title ?? settings.meta?.title ?? settings.siteName ?? siteName,
		description: pageMeta?.description ?? settings.meta?.description ?? settings.tagLine ?? '',
		media: pageMeta?.image
			? resolveItems(pageMeta.image)
			: settings.meta?.image
				? resolveItems(settings.meta.image)
				: null,
		canonical: pageMeta?.canonicalURL ?? settings.meta?.canonicalURL ?? null,
		private: pageMeta?.private ?? settings.meta?.private ?? false,
		structuredData: {
			settings: settings.meta?.structuredData ?? null,
			page: pageMeta?.structuredData ?? null,
		},
	};
});

/**
 * Other meta props that could change on page refresh.
 */
const site = $derived(settings.siteName?.trim() ?? siteName);
const locale = $derived(settings?.locale ?? 'en');
const social = $derived(settings?.social ?? {});
const codeInjectionSettings = $derived(settings?.codeInjection);
const socialLinks = $derived(Object.values(social).filter(Boolean) as string[]);

/**
 * Generate structured data for the website.
 */
const websiteSchema = $derived.by(() => ({
	'@context': 'https://schema.org',
	'@type': 'WebSite',
	name: siteName,
	url: meta.canonical,
}));

/**
 * Generate structured data for the organisation.
 */
const organizationSchema = $derived.by(() => {
	if (socialLinks.length === 0) return null;

	return {
		'@context': 'https://schema.org',
		'@type': 'Organization',
		name: siteName,
		...(meta.canonical ? { url: meta.canonical } : {}),
		...(meta?.media?.url ? { logo: meta.media.url } : {}),
		...(settings.tagLine || settings.meta?.description
			? { description: settings.tagLine ?? settings.meta?.description }
			: {}),
		sameAs: socialLinks.filter((url) => url.startsWith('http')),
		...(settings.contact?.email || settings.contact?.telephone
			? {
					contactPoint: [
						{
							'@type': 'ContactPoint',
							...(settings.contact?.telephone
								? { telephone: settings.contact.telephone }
								: {}),
							...(settings.contact?.email ? { email: settings.contact.email } : {}),
							contactType: 'customer support',
						},
					],
				}
			: {}),
		...(settings.address
			? {
					address: {
						'@type': 'PostalAddress',
						...(settings.address.line1
							? { streetAddress: settings.address.line1 }
							: {}),
						...(settings.address.city
							? { addressLocality: settings.address.city }
							: {}),
						...(settings.address.county
							? { addressRegion: settings.address.county }
							: {}),
						...(settings.address.postcode
							? { postalCode: settings.address.postcode }
							: {}),
						...(settings.address.country
							? { addressCountry: settings.address.country }
							: {}),
					},
				}
			: {}),
	};
});
</script>

<!--
	@component

	PayloadSEO renders all head meta tags for a page, including Open Graph,
	Twitter Card, canonical URL, structured data, and code injection from
	Payload CMS settings.

	@example
	```svelte
	<PayloadSEO
		siteName="My Site"
		settings={globalSettings}
		pageMeta={page.meta}
	/>
	```
-->
<svelte:head>
	<!-- Meta -->
	<title>{meta.title} | {siteName}</title>
	<meta name="description" content={meta.description} />
	<!-- Canonical -->
	{#if meta.canonical}
		<link rel="canonical" href={meta.canonical} />
	{/if}
	<!-- No Index -->
	{#if meta.private}
		<meta name="robots" content="noindex" />
	{/if}
	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:title" content={meta.title} />
	<meta property="og:description" content={meta.description} />
	<meta property="og:site_name" content={site} />
	<meta property="og:locale" content={locale} />
	<meta property="og:url" content={meta.canonical} />
	{#if meta.media}
		{#if meta.media.url}
			<meta property="og:image" content={meta.media.url} />
		{/if}
		{#if meta.media.alt}
			<meta property="og:image:alt" content={meta.media.alt} />
		{/if}
		{#if meta.media.width}
			<meta property="og:image:width" content={meta.media.width.toString()} />
		{/if}
		{#if meta.media.height}
			<meta property="og:image:height" content={meta.media.height.toString()} />
		{/if}
		{#if meta.media.mimeType}
			<meta property="og:image:type" content={meta.media.mimeType} />
		{/if}
	{/if}
	<!-- Twitter Card -->
	<meta name="twitter:card" content={meta?.media?.url ? 'summary_large_image' : 'summary'} />
	<meta name="twitter:title" content={meta.title} />
	<meta name="twitter:description" content={meta.description} />
	{#if meta.media}
		{#if meta.media.url}
			<meta name="twitter:image" content={meta.media.url} />
		{/if}
		{#if meta.media.alt}
			<meta name="twitter:image:alt" content={meta.media.alt} />
		{/if}
	{/if}
	{#if social.x}
		<meta name="twitter:site" content={social.x.startsWith('@') ? social.x : `@${social.x}`} />
	{/if}
	<!-- Structured Data - Settings -->
	{#if meta.structuredData.settings}
		{@html serializeSchema(meta.structuredData.settings)}
	{/if}
	<!-- Structured Data - Page -->
	{#if meta.structuredData.page}
		{@html serializeSchema(meta.structuredData.page)}
	{/if}
	<!-- Website Schema -->
	{#if websiteSchema}
		{@html serializeSchema(websiteSchema)}
	{/if}
	<!-- Social Organisation Schema -->
	{#if organizationSchema}
		{@html serializeSchema(organizationSchema)}
	{/if}
	<!-- Code Injection - Settings -->
	{#if codeInjectionSettings?.head}
		{@html codeInjectionSettings.head}
	{/if}
</svelte:head>
