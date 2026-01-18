<script lang="ts" module>
import type { Snippet } from 'svelte';

export type SidebarProps = {
	menuLabel?: string;
	children: Snippet;
	isOpen?: boolean;
	position?: 'left' | 'right';
	width?: string;
	top?: number;
	closeOnOverlayClick?: boolean;
	overlayOpacity?: number;
	toggleStyle?: 'toggle' | 'hamburger';
	class?: string;
	onOpen?: () => void;
	onClose?: () => void;
	onToggle?: (isOpen: boolean) => void;
};
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import Hamburger from './Hamburger.svelte';

	let {
		menuLabel = 'Menu',
		children,
		isOpen = $bindable(false),
		position = 'left',
		width = '50vw',
		top = 160,
		closeOnOverlayClick = true,
		overlayOpacity = 0.3,
		toggleStyle = 'toggle',
		class: className = '',
		onOpen,
		onClose,
		onToggle
	}: SidebarProps = $props();

	// Sync checkbox with isOpen state
	$effect(() => {
		const checkbox = document.querySelector<HTMLInputElement>('.sidebar__checkbox');
		if (checkbox && checkbox.checked !== isOpen) {
			checkbox.checked = isOpen;
		}
	});

	// Watch for changes to isOpen and call callbacks
	$effect(() => {
		if (isOpen) {
			onOpen?.();
		} else {
			onClose?.();
		}
		onToggle?.(isOpen);
	});

	onMount(() => {
		const overlay = document.querySelector('.sidebar__overlay');
		const checkbox = document.querySelector<HTMLInputElement>('.sidebar__checkbox');

		if (overlay && checkbox) {
			const handleOverlayClick = (e: Event) => {
				if (!closeOnOverlayClick) return;
				e.preventDefault();
				checkbox.checked = false;
				isOpen = false;
			};

			const handleCheckboxChange = () => {
				isOpen = checkbox.checked;
			};

			overlay.addEventListener('click', handleOverlayClick);
			checkbox.addEventListener('change', handleCheckboxChange);

			return () => {
				overlay.removeEventListener('click', handleOverlayClick);
				checkbox.removeEventListener('change', handleCheckboxChange);
			};
		}
	});
</script>

<!--
	@component

	Mobile-first sidebar navigation component with toggle and hamburger options.
	Automatically collapses on mobile and remains visible on desktop.

	@example
	```svelte
	<Sidebar bind:isOpen>
		<nav>
			<a href="/">Home</a>
			<a href="/about">About</a>
		</nav>
	</Sidebar>
	```

	@example
	```svelte
	<Sidebar toggleStyle="hamburger" position="right" width="300px">
		<nav>...</nav>
	</Sidebar>
	```
-->
<aside
	class="sidebar sidebar--{toggleStyle} sidebar--{position} {className}"
	style="--sidebar-width: {width}; --sidebar-top: {top}px; --sidebar-overlay-opacity: {overlayOpacity}"
>
	<!-- Click Logic -->
	<input type="checkbox" class="sidebar__checkbox" id="sidebar-checkbox" />
	<label for="sidebar-checkbox" class="sidebar__overlay"></label>

	<!-- Hamburger Toggle -->
	{#if toggleStyle === 'hamburger'}
		<Hamburger bind:isOpen />
	{/if}

	<!-- Content -->
	<div class="sidebar__content">
		{#if toggleStyle === 'toggle'}
			<label for="sidebar-checkbox" class="sidebar__toggle">
				{menuLabel}
			</label>
		{/if}
		<div class="sidebar__inner">
			{@render children()}
		</div>
	</div>
</aside>

<style lang="scss">
	.sidebar {
		$self: &;

		&__toggle {
			position: absolute;
			display: none;
			bottom: 0;
			right: 1px;
			background-color: var(--sidebar-toggle-background, var(--colour-base-black));
			color: var(--sidebar-toggle-colour, var(--colour-base-light));
			padding: var(--sidebar-toggle-padding, 0.25rem 1.5rem);
			border-top-right-radius: var(--sidebar-toggle-radius, 0.375rem);
			border-top-left-radius: var(--sidebar-toggle-radius, 0.375rem);
			font-size: var(--sidebar-toggle-font-size, 0.9rem);
			transform: rotate(90deg) translate(0%, -100%);
			transform-origin: right top;
			cursor: pointer;
			user-select: none;
			transition: box-shadow 200ms linear;
			border: 1px solid var(--sidebar-border-colour, rgba(255, 255, 255, 0.1));

			&::before {
				content: '';
				position: absolute;
				top: calc(90% + 2px);
				left: 1px;
				width: calc(100% - 2px);
				height: 10%;
				background: var(--sidebar-toggle-background, var(--colour-base-black));
			}
		}

		&__overlay {
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: var(--sidebar-overlay-colour, var(--colour-grey-900));
			z-index: -100;
			opacity: 0;
			transition:
				opacity 400ms ease,
				z-index 400ms step-end;
		}

		&__checkbox {
			position: fixed;
			top: 0;
			display: none;

			&:checked {
				~ #{$self}__content {
					translate: 0;
					z-index: 9999999;

					#{$self}__toggle {
						box-shadow: none;
					}
				}

				~ #{$self}__overlay {
					transition:
						opacity 600ms ease,
						z-index 600ms step-start;
					opacity: var(--sidebar-overlay-opacity, 0.3);
					z-index: 999999;
				}
			}
		}

		@media (max-width: 1023px) {
			&__content {
				position: fixed;
				display: grid;
				top: 0;
				height: 100%;
				width: var(--sidebar-width, 50vw);
				min-width: var(--sidebar-min-width, 270px);
				background-color: var(--sidebar-background, var(--colour-base-black));
				border-color: var(--sidebar-border-colour, rgba(255, 255, 255, 0.1));
				z-index: 1000;
				transition: translate 600ms cubic-bezier(0.1, 0.7, 0.1, 1);
			}

			&__inner {
				overflow: auto;
				display: flex;
				flex-direction: column;
				padding: var(--sidebar-inner-padding, 2rem 1.8rem 0 1.8rem);
			}

			&__toggle {
				display: flex;
			}
		}

		@media (min-width: 1024px) {
			position: sticky;
			top: var(--sidebar-top, 160px);

			&__overlay {
				display: none;
			}
		}

		&--left {
			@media (max-width: 1023px) {
				#{$self}__content {
					left: 0;
					border-right-style: solid;
					border-right-width: 1px;
					translate: -100%;
				}
			}
		}

		&--right {
			@media (max-width: 1023px) {
				#{$self}__content {
					right: 0;
					border-left-style: solid;
					border-left-width: 1px;
					translate: 100%;
				}
			}
		}

		&--hamburger {
			#{$self}__toggle {
				display: none;
			}
		}
	}
</style>
