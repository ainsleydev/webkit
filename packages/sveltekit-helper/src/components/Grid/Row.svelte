<script lang="ts" module>
export type RowProps = {
	noGaps?: boolean;
};
</script>

<script lang="ts">
const { noGaps = false, ...restProps }: RowProps = $props();
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
<div class="row" class:row--no-gaps={noGaps} {...restProps}>
	<slot />
</div>

<style lang="scss">
	.row {
		$self: &;

		display: flex;
		flex-wrap: wrap;
		margin-inline: calc(var(--row-gap, 1rem) * -1);

		&--no-gaps {
			margin-inline: 0;

			:global(.col),
			:global([class*='col-']) {
				padding-inline: 0;
			}
		}

		@media (max-width: 568px) {
			margin-inline: calc(var(--row-gap-mobile, var(--row-gap, 0.5rem)) * -1);
		}
	}
</style>
