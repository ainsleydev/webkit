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
-->
<div class="col" {...restProps}>
	<slot />
</div>

<style lang="scss">
	.col {
		--col-gap: var(--col-gap, 1rem);

		position: relative;
		width: 100%;
		padding-inline: var(--col-gap);

		@media (max-width: 568px) {
			--col-gap: var(--col-gap-mobile, var(--col-gap, 0.5rem));
		}
	}
</style>
