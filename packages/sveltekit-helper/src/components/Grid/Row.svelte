<script lang="ts">
let noGaps = false;
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

<style lang="scss">
	.row {
		$self: &;

		--row-gap: 1rem;

		display: flex;
		flex-wrap: wrap;
		margin-inline: calc(var(--row-gap) * -1);

		&--no-gaps {
			margin-inline: 0;

			:global(.col),
			:global([class*='col-']) {
				padding-inline: 0;
			}
		}

		@media (max-width: 568px) {
			--row-gap: 0.5rem;
		}
	}
</style>
