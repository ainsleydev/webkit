<script lang="ts">
import { serializeSchema } from '../../utils/seo/ld-json.js';

import type { PayloadFooterProps } from './types.js';

const data: PayloadFooterProps = $props();

/**
 * Code injection meta.
 */
const siteCodeInjection = $derived(data.settings?.codeInjection?.footer);
const pageCodeInjection = $derived(data.pageCodeInjection?.footer);
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
{#if pageCodeInjection}
	{@html pageCodeInjection}
{/if}
