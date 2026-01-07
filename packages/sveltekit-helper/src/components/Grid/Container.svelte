<script lang="ts">
let props: Record<string, unknown> = $props();
</script>

<!--
	@component

	Centre content such as rows & columns horizontally with predefined max-width.
	Uses CSS Grid to provide breakout and full-width layout options.

	@example
	```svelte
	<Container>
		<Row></Row>
	</Container>
	```

	@example
	```svelte
	<!-- Custom max width via CSS variable -->
	<Container style="--container-max-width: 1400px">
		<Row></Row>
	</Container>
	```
-->
<div class="container" {...props}>
	<slot />
</div>

<style lang="scss">
	.container {
		$self: &;

		--container-padding: 1rem;
		--container-max-width: 1328px;
		--container-breakout-max-width: 1500px;
		--container-breakout-size: calc(
			(var(--container-breakout-max-width) - var(--container-max-width)) / 2
		);

		display: grid;
		width: 100%;
		position: relative;
		grid-template-columns:
			[full-width-start] minmax(var(--container-padding), 1fr)
			[breakout-start] minmax(0, var(--container-breakout-size))
			[content-start] min(
				100% - (var(--container-padding) * 2),
				var(--container-max-width)
			)
			[content-end]
			minmax(0, var(--container-breakout-size)) [breakout-end]
			minmax(var(--container-padding), 1fr) [full-width-end];

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
