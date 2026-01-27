<script lang="ts">
const { ...restProps } = $props();
</script>

<!--
	@component

	Responsive column component.
	Provides base column structure with customisable gap.
	Consumers should define their own grid classes (col-12, col-tab-6, etc.) in their global styles.

	@example
	```svelte
	<Column class="col-12 col-tab-6 col-desk-4">
		Content
	</Column>
	```

	CSS Custom Properties:
	- `--col-gap`: Horizontal padding on desktop (default: 1rem)
	- `--col-gap-mobile`: Horizontal padding on mobile (default: var(--col-gap, 0.5rem))
-->
<div class="col" {...restProps}>
	<slot />
</div>

<style lang="scss">
	.col {
		position: relative;
		width: 100%;
		padding-inline: var(--col-gap, 1rem);

		@media (max-width: 568px) {
			padding-inline: var(--col-gap-mobile, var(--col-gap, 0.5rem));
		}
	}
</style>
