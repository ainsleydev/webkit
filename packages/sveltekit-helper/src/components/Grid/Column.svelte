<!--
	@component

	Responsive column component with 12-column grid system.
	Columns build on the responsive grid system by allowing how columns
	grow, shrink and change at different viewports.

	@example
	```svelte
	<Column class="col-12 col-tab-10 col-desk-4">
		Content
	</Column>
	```
-->
<div class="col-12" {...$$restProps}>
	<slot />
</div>

<style>
	:global([class*='col-']) {
		--col-gap: 1rem;

		position: relative;
		width: 100%;
		padding-inline: var(--col-gap);
	}

	/* Base columns */
	:global(.col-auto) {
		flex: 0 0 auto;
		width: auto;
		max-width: 100%;
	}

	:global(.col-1) { width: 8.333333%; }
	:global(.col-2) { width: 16.666667%; }
	:global(.col-3) { width: 25%; }
	:global(.col-4) { width: 33.333333%; }
	:global(.col-5) { width: 41.666667%; }
	:global(.col-6) { width: 50%; }
	:global(.col-7) { width: 58.333333%; }
	:global(.col-8) { width: 66.666667%; }
	:global(.col-9) { width: 75%; }
	:global(.col-10) { width: 83.333333%; }
	:global(.col-11) { width: 91.666667%; }
	:global(.col-12) { width: 100%; }

	/* Tablet breakpoint (768px) */
	@media (min-width: 768px) {
		:global(.col-tab-auto) {
			flex: 0 0 auto;
			width: auto;
			max-width: 100%;
		}

		:global(.col-tab-1) { width: 8.333333%; }
		:global(.col-tab-2) { width: 16.666667%; }
		:global(.col-tab-3) { width: 25%; }
		:global(.col-tab-4) { width: 33.333333%; }
		:global(.col-tab-5) { width: 41.666667%; }
		:global(.col-tab-6) { width: 50%; }
		:global(.col-tab-7) { width: 58.333333%; }
		:global(.col-tab-8) { width: 66.666667%; }
		:global(.col-tab-9) { width: 75%; }
		:global(.col-tab-10) { width: 83.333333%; }
		:global(.col-tab-11) { width: 91.666667%; }
		:global(.col-tab-12) { width: 100%; }

		:global(.offset-tab-0) { margin-left: 0; }
		:global(.offset-tab-1) { margin-left: 8.333333%; }
		:global(.offset-tab-2) { margin-left: 16.666667%; }
		:global(.offset-tab-3) { margin-left: 25%; }
		:global(.offset-tab-4) { margin-left: 33.333333%; }
		:global(.offset-tab-5) { margin-left: 41.666667%; }
		:global(.offset-tab-6) { margin-left: 50%; }
		:global(.offset-tab-7) { margin-left: 58.333333%; }
		:global(.offset-tab-8) { margin-left: 66.666667%; }
		:global(.offset-tab-9) { margin-left: 75%; }
		:global(.offset-tab-10) { margin-left: 83.333333%; }
		:global(.offset-tab-11) { margin-left: 91.666667%; }
	}

	/* Desktop breakpoint (1024px) */
	@media (min-width: 1024px) {
		:global(.col-desk-auto) {
			flex: 0 0 auto;
			width: auto;
			max-width: 100%;
		}

		:global(.col-desk-1) { width: 8.333333%; }
		:global(.col-desk-2) { width: 16.666667%; }
		:global(.col-desk-3) { width: 25%; }
		:global(.col-desk-4) { width: 33.333333%; }
		:global(.col-desk-5) { width: 41.666667%; }
		:global(.col-desk-6) { width: 50%; }
		:global(.col-desk-7) { width: 58.333333%; }
		:global(.col-desk-8) { width: 66.666667%; }
		:global(.col-desk-9) { width: 75%; }
		:global(.col-desk-10) { width: 83.333333%; }
		:global(.col-desk-11) { width: 91.666667%; }
		:global(.col-desk-12) { width: 100%; }

		:global(.offset-desk-0) { margin-left: 0; }
		:global(.offset-desk-1) { margin-left: 8.333333%; }
		:global(.offset-desk-2) { margin-left: 16.666667%; }
		:global(.offset-desk-3) { margin-left: 25%; }
		:global(.offset-desk-4) { margin-left: 33.333333%; }
		:global(.offset-desk-5) { margin-left: 41.666667%; }
		:global(.offset-desk-6) { margin-left: 50%; }
		:global(.offset-desk-7) { margin-left: 58.333333%; }
		:global(.offset-desk-8) { margin-left: 66.666667%; }
		:global(.offset-desk-9) { margin-left: 75%; }
		:global(.offset-desk-10) { margin-left: 83.333333%; }
		:global(.offset-desk-11) { margin-left: 91.666667%; }
	}

	@media (max-width: 568px) {
		:global([class*='col-']) {
			--col-gap: 0.5rem;
		}
	}
</style>
