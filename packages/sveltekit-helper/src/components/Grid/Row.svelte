<script lang="ts">
export const noGaps = false;
</script>

<!--
	@component

	Flexbox row container for columns with gap management.

	@example
	```svelte
	<Row>
		<Column></Column>
	</Row>
	```

	@example
	```svelte
	<Row noGaps>
		<Column></Column>
	</Row>
	```
-->
<div class="row" class:row--no-gaps={noGaps} {...$$restProps}>
	<slot />
</div>

<style>
	.row {
		--row-gap: 1rem;

		display: flex;
		flex-wrap: wrap;
		margin-inline: calc(var(--row-gap) * -1);
	}

	.row--no-gaps {
		margin-inline: 0;
	}

	.row--no-gaps :global(.col),
	.row--no-gaps :global([class*='col-']) {
		padding-inline: 0;
	}

	@media (max-width: 568px) {
		.row {
			--row-gap: 0.5rem;
		}
	}
</style>
