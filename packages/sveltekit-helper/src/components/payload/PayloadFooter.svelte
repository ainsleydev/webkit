<script lang="ts">
import type { PayloadFooterProps } from './types.js';

const { settings, pageCodeInjection }: PayloadFooterProps = $props();

/**
 * Code injection meta.
 */
const siteCodeInjection = $derived(settings?.codeInjection?.footer);
const pageFooterInjection = $derived(pageCodeInjection?.footer);
</script>

<!--
	@component

	PayloadFooter renders footer code injection scripts from both
	site-level settings and page-level overrides.

	@example
	```svelte
	<PayloadFooter
		settings={globalSettings}
		pageCodeInjection={page.codeInjection}
	/>
	```
-->
<!-- Code Injection - Settings -->
{#if siteCodeInjection}
	{@html siteCodeInjection}
{/if}
<!-- Code Injection - Page -->
{#if pageFooterInjection}
	{@html pageFooterInjection}
{/if}
