<script lang="ts">
const { ...restProps } = $props();
</script>

<!--
	@component

	Centre content such as rows & columns horizontally with predefined max-width.
	Uses CSS Grid to provide breakout and full-width layout options.

	@example
	```svelte
	<Container>
		<Row>
			<Column class="col-12">Content</Column>
		</Row>
	</Container>
	```

	@example
	```svelte
	<Container>
		<div class="breakout">Full breakout content</div>
		<div class="full-width">Full width content</div>
	</Container>
	```

	CSS Custom Properties:
	- `--container-max-width`: Maximum content width (default: 1328px)
	- `--container-breakout-max-width`: Maximum breakout width (default: 1500px)
	- `--container-padding`: Horizontal padding (default: 1rem)
-->
<div class="container" {...restProps}>
	<slot />
</div>

<style lang="scss">
	.container {
		$self: &;

		--container-breakout-size: calc(
			(var(--container-breakout-max-width, 1500px) - var(--container-max-width, 1328px)) / 2
		);

		display: grid;
		width: 100%;
		position: relative;
		grid-template-columns:
			[full-width-start] minmax(var(--container-padding, 1rem), 1fr)
			[breakout-start] minmax(0, var(--container-breakout-size))
			[content-start] min(
				100% - (var(--container-padding, 1rem) * 2),
				var(--container-max-width, 1328px)
			)
			[content-end]
			minmax(0, var(--container-breakout-size)) [breakout-end]
			minmax(var(--container-padding, 1rem), 1fr) [full-width-end];

		:global(> *) {
			grid-column: content;
		}

		:global(> .breakout) {
			grid-column: breakout;
		}

		:global(> .full-width) {
			display: grid;
			grid-column: full-width;
			grid-template-columns: inherit;
		}
	}
</style>
