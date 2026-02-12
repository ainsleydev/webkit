<script lang="ts" module>
import type { Snippet, TransitionConfig } from 'svelte';

export type TransitionFn = (node: Element, params: Record<string, unknown>) => TransitionConfig;

export type ModalProps = {
	title?: string;
	isOpen?: boolean;
	children?: Snippet;
	class?: string;
	onClose?: () => void;
	transition?: TransitionFn;
	transitionParams?: Record<string, unknown>;
};
</script>

<script lang="ts">
	import { X } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';

	let {
		title = '',
		isOpen = $bindable(false),
		children,
		class: className = '',
		onClose,
		transition: transitionFn = fade,
		transitionParams = { duration: 100 }
	}: ModalProps = $props();

	let modalContent = $state<HTMLDivElement>();

	const handleBackdropClick = (event: MouseEvent) => {
		if (modalContent && !modalContent.contains(event.target as Node)) {
			onClose?.();
		}
	};

	onMount(() => {
		const handleKeydown = (e: KeyboardEvent) => {
			if (e.key === 'Escape' && isOpen) {
				e.preventDefault();
				onClose?.();
			}
		};

		document.addEventListener('keydown', handleKeydown);

		return () => {
			document.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

<!--
	@component

	Modal dialog component with backdrop click and Escape key close behaviour.
	Uses the native `<dialog>` element for accessibility.

	@example Basic usage with title
	```svelte
	<Modal bind:isOpen title="Confirm action" onClose={() => (isOpen = false)}>
		<p>Are you sure you want to proceed?</p>
	</Modal>
	```

	@example Without title
	```svelte
	<Modal bind:isOpen onClose={() => (isOpen = false)}>
		<form>...</form>
	</Modal>
	```

	@example Slide in from the left
	```svelte
	<script>
		import { fly } from 'svelte/transition';
	</script>

	<Modal
		bind:isOpen
		title="Slide modal"
		onClose={() => (isOpen = false)}
		transition={fly}
		transitionParams={{ x: -300, duration: 200 }}
	/>
	```

	@example Scale up from centre
	```svelte
	<script>
		import { scale } from 'svelte/transition';
	</script>

	<Modal
		bind:isOpen
		onClose={() => (isOpen = false)}
		transition={scale}
		transitionParams={{ start: 0.9, duration: 150 }}
	>
		<p>Scaled content</p>
	</Modal>
	```

	CSS Custom Properties:
	- `--modal-overlay-bg`: Backdrop colour (default: rgba(0, 0, 0, 0.6))
	- `--modal-padding-top`: Top offset from viewport (default: var(--header-height))
	- `--modal-content-max-width`: Max width of the content panel (default: 600px)
	- `--modal-content-bg`: Content background (default: var(--token-surface-default))
	- `--modal-content-border`: Content border (default: 1px solid var(--token-border-grey))
	- `--modal-content-border-radius`: Content border radius (default: 12px)
	- `--modal-content-padding`: Content padding (default: 1.5rem / 2rem on tablet)
	- `--modal-header-border`: Header bottom border (default: 1px solid var(--token-border-grey))
	- `--modal-close-colour`: Close icon colour (default: var(--token-icon-grey))
-->
{#if isOpen}
	<dialog
		open
		class="modal {className}"
		aria-modal="true"
		aria-label={title || undefined}
		onclick={handleBackdropClick}
		transition:transitionFn={transitionParams}
	>
		<div class="modal__content" bind:this={modalContent}>
			{#if title}
				<header class="modal__header">
					<h4 class="modal__title">{title}</h4>
					<button
						class="modal__close"
						onclick={() => onClose?.()}
						aria-label={title ? `Close ${title}` : 'Close modal'}
					>
						<X color="var(--modal-close-colour, var(--token-icon-grey))" />
					</button>
				</header>
			{/if}
			<div class="modal__body">
				{#if children}
					{@render children()}
				{/if}
			</div>
		</div>
	</dialog>
{/if}

<style lang="scss">
	@use '../scss' as a;

	.modal {
		$self: &;

		position: fixed;
		display: flex;
		align-items: flex-start;
		justify-content: center;
		top: 0;
		left: 0;
		padding: var(--modal-padding-top, var(--header-height)) 0 0;
		height: 100vh;
		width: 100vw;
		max-width: none;
		max-height: none;
		background-color: var(--modal-overlay-bg, rgba(0, 0, 0, 0.6));
		z-index: 9999999;
		outline: none;
		border: none;

		&__header {
			display: flex;
			justify-content: space-between;
			align-items: flex-start;
			width: 100%;
			border-bottom: var(--modal-header-border, 1px solid var(--token-border-grey));
			margin-bottom: a.$size-16;
			padding-bottom: a.$size-16;
		}

		&__title {
			margin: 0;
		}

		&__close {
			cursor: pointer;
			background: none;
			border: none;
			padding: 0;
			display: flex;
			align-items: center;
			justify-content: center;
		}

		&__content {
			position: relative;
			display: flex;
			flex-direction: column;
			align-items: flex-start;
			width: min(calc(100% - 1.6rem), var(--modal-content-max-width, 600px));
			background: var(--modal-content-bg, var(--token-surface-default));
			padding: var(--modal-content-padding, a.$size-24);
			border: var(--modal-content-border, 1px solid var(--token-border-grey));
			border-radius: var(--modal-content-border-radius, a.$border-radius-12);
		}

		@include a.mq(tablet) {
			&__content {
				padding: var(--modal-content-padding, a.$size-32);
			}
		}
	}
</style>
