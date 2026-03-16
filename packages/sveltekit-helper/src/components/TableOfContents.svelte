<script lang="ts" module>
export type TOCItem = {
	label: string;
	href: string;
};

export type TableOfContentsProps = {
	/**
	 * Optional heading displayed above the TOC list.
	 */
	heading?: string;
	/**
	 * Optional pre-generated items. If omitted, items are auto-discovered
	 * from the DOM using contentSelector + headingSelector on mount.
	 */
	items?: TOCItem[];
	/**
	 * Adds a left border to the TOC on tablet and above.
	 */
	displayBorder?: boolean;
	/**
	 * CSS selector used to find the content element containing headings.
	 * Defaults to `[data-sidebar-content="true"]`.
	 */
	contentSelector?: string;
	/**
	 * CSS selector for headings within the content element.
	 * Falls back to the `data-sidebar-selector` attribute on the content
	 * element, then `h2`.
	 *
	 * Priority: prop > data-sidebar-selector attribute > 'h2'
	 */
	headingSelector?: string;
	/**
	 * Scroll offset in pixels applied to scrollspy detection.
	 * @default 80
	 */
	scrollOffset?: number;
	/**
	 * Optional callback invoked when a TOC link is clicked.
	 * Useful for closing a sidebar or drawer on navigation.
	 */
	onLinkClick?: (event: MouseEvent, item: TOCItem) => void;
};
</script>

<script lang="ts">
	import { onMount } from 'svelte'

	let {
		heading = '',
		items: itemsProp,
		displayBorder = false,
		contentSelector = '[data-sidebar-content="true"]',
		headingSelector,
		scrollOffset = 80,
		onLinkClick,
	}: TableOfContentsProps = $props()

	let activeId = $state<string | null>(null)
	let items = $state<TOCItem[]>(itemsProp ?? [])

	onMount(() => {
		const content = document.querySelector(contentSelector)
		if (!content) return

		// Priority: prop > data-sidebar-selector attribute > 'h2'
		const resolvedHeadingSelector =
			headingSelector ?? content.getAttribute('data-sidebar-selector') ?? 'h2'

		// Auto-generate items from DOM if not provided as props.
		if (!itemsProp) {
			items = Array.from(content.querySelectorAll<HTMLElement>(resolvedHeadingSelector))
				.filter((el) => el.id)
				.map((el) => ({ label: el.textContent?.trim() ?? '', href: el.id }))
		}

		const sections = Array.from(
			content.querySelectorAll<HTMLElement>(resolvedHeadingSelector),
		).filter((el) => el.id)

		if (sections.length === 0) return

		const onScroll = () => {
			const scrollPosition = window.scrollY + scrollOffset
			let activeSection = sections[0]

			for (const section of sections) {
				if (section.offsetTop > scrollPosition) break
				activeSection = section
			}

			const bottomThreshold = 10
			const scrolledToBottom =
				window.innerHeight + window.scrollY >= document.body.scrollHeight - bottomThreshold
			if (scrolledToBottom) {
				activeSection = sections[sections.length - 1]
			}

			activeId = activeSection?.id ?? null
		}

		window.addEventListener('scroll', onScroll, { passive: true })
		onScroll()

		return () => window.removeEventListener('scroll', onScroll)
	})
</script>

<!--
	@component

	Table of Contents with scrollspy, designed to be used alongside a richtext
	or content area that has headings with `id` attributes.

	By default the component discovers the content element via a
	`data-sidebar-content="true"` attribute and uses the `data-sidebar-selector`
	attribute (defaulting to `h2`) to determine which headings to track.

	@example
	```svelte
	<RichText content={data.body} data-sidebar-content="true" data-sidebar-selector="h3" />
	<TableOfContents heading="On this page" />
	```

	@example
	```svelte
	<TableOfContents
		contentSelector=".article-body"
		headingSelector="h2, h3"
		heading="Contents"
		displayBorder
	/>
	```

	@example
	```svelte
	<TableOfContents items={[{ label: 'Intro', href: 'intro' }]} />
	```
-->
<div class="toc" class:toc--border={displayBorder}>
	{#if heading !== ''}
		<p class="toc__heading">
			{heading}
		</p>
	{/if}
	<menu class="toc__items">
		{#each items as item, index (index)}
			<li class="toc__item">
				<a
					class="toc__link"
					class:toc__link--active={activeId ===
						(item.href.startsWith('#') ? item.href.slice(1) : item.href)}
					href="#{item.href}"
					onclick={(e) => onLinkClick?.(e, item)}
				>
					<small>{item.label}</small>
				</a>
			</li>
		{/each}
	</menu>
</div>

<style lang="scss">
	@use '../scss' as a;

	.toc {
		$self: &;

		&__item {
			margin-bottom: a.$size-8;
		}

		&__link {
			text-decoration: none;
			font-weight: var(--font-weight-normal);
			color: var(--token-colour-text);
			transition: color 100ms ease;
			will-change: color;

			&--active {
				color: var(--toc-colour-active, var(--token-text-action));
				font-weight: var(--font-weight-medium);
			}

			&:hover {
				color: var(--toc-colour-active, var(--token-text-action));
			}
		}

		@include a.mq(tablet) {
			&--border {
				margin-left: var(--toc-border-offset, #{a.$size-48});
				padding-left: var(--toc-border-offset, #{a.$size-48});
				border-left: 1px solid var(--toc-border-colour, var(--colour-light-600));
			}
		}
	}
</style>
