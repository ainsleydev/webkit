<script lang="ts" module>
	import type { HamburgerType as HamburgerTypeLib } from "svelte-hamburgers";

	export type HamburgerType = HamburgerTypeLib

	export type HamburgerProps = {
		// See: https://github.com/ghostdevv/svelte-hamburgers/blob/main/types.md
		style?: HamburgerType;
		isOpen?: boolean;
		gap?: string;
		class?: string;
		onChange?: (isOpen: boolean) => void;
	};
</script>

<script lang="ts">
	import { Hamburger as SvelteHamburger } from 'svelte-hamburgers';

	let {
		style = 'spin',
		isOpen = $bindable(false),
		gap = '0.8rem',
		class: className = '',
		onChange
	}: HamburgerProps = $props();
</script>

<!--
	@component

	Hamburger menu icon with animation for mobile navigation.
	Uses svelte-hamburgers under the hood.

	@example
	```svelte
	<Hamburger bind:isOpen />
	```

	@example
	```svelte
	<Hamburger gap="1rem">
	```
-->
<div class="hamburger-wrapper {className}" style="--hamburger-gap: {gap}" aria-label="Toggle Menu">
	<SvelteHamburger
		type={style}
		bind:open={isOpen}
		on:change={() => onChange?.(isOpen)}
		--color="var(--hamburger-colour, var(--colour-base-light))"
		--layer-width="var(--hamburger-layer-width, 24px)"
		--layer-height="var(--hamburger-layer-height, 2px)"
		--layer-spacing="var(--hamburger-layer-spacing, 5px)"
		--border-radius="var(--hamburger-border-radius, 2px)"
	/>
</div>

<style lang="scss">
	.hamburger-wrapper {
		position: fixed;
		top: var(--hamburger-gap, 0.8rem);
		right: var(--hamburger-gap, 0.8rem);
		z-index: var(--hamburger-z-index, 10000);
	}
</style>
