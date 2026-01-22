<script lang="ts" module>
import type { Icon as IconType } from '@lucide/svelte';

export type NoticeType = 'info' | 'warning' | 'success' | 'error';

export type NoticeProps = {
	type?: NoticeType;
	title: string;
	visible?: boolean;
	dismiss?: boolean;
	icon?: typeof IconType;
};
</script>

<script lang="ts">
	import { X } from '@lucide/svelte'
	import { fade } from 'svelte/transition'

	import { alertIcons } from './alertIcons'

	let {
		type = 'info',
		title = '',
		visible = $bindable(true),
		dismiss = false,
		icon: customIcon,
		...restProps
	}: NoticeProps = $props()

	const iconDetail = $derived(alertIcons[type])
	const Icon = $derived(customIcon || iconDetail.icon)
	const hide = () => (visible = false)
	const ariaLive = $derived(type === 'error' ? 'assertive' : 'polite')
</script>

<!--
	@component

	Inline notification component for displaying brief messages with icons.
	Compact design suitable for inline alerts, badges, or status indicators.

	@example
	```svelte
	<Notice type="success" title="Upload complete" />
	<Notice type="warning" title="Session expiring" dismiss />
	<Notice type="error" title="Connection failed" icon={CustomIcon} />
	```

	CSS Custom Properties:
	- `--_notice-gap`: Gap between icon and title (default: 12px)
	- `--_notice-padding`: Internal padding (default: 0.8rem 12px)
	- `--_notice-border-radius`: Border radius (default: 6px)
	- `--_notice-bg`: Background color (default: rgba(255, 255, 255, 0.025))
	- `--_notice-font-size`: Title font size (default: 1rem)
	- `--_notice-title-colour`: Title text color (default: rgba(255, 255, 255, 1))
	- `--_notice-icon-colour`: Icon color (set automatically based on type)
-->
{#if visible}
	<div
		class="notice notice--{type}"
		role="alert"
		aria-live={ariaLive}
		aria-atomic="true"
		style="--_notice-icon-colour: {iconDetail.colour}"
		transition:fade={{ duration: 300 }}
		{...restProps}
	>
		<!-- Icon -->
		<figure class="notice__icon">
			<Icon size={20} color={iconDetail.colour} strokeWidth={1.2}></Icon>
		</figure>
		<!-- Title -->
		<p class="notice__title">
			{title}
		</p>
		<!-- Dismiss -->
		{#if dismiss}
			<button class="notice__dismiss" onclick={hide} aria-label="Close Notice">
				<X size={20} color={iconDetail.colour} />
			</button>
		{/if}
	</div>
{/if}

<style lang="scss">
	@use '../../scss' as a;

	.notice {
		$self: &;

		position: relative;
		display: inline-flex;
		width: auto;
		gap: var(--_notice-gap, a.$size-12);
		padding: var(--_notice-padding, 0.8rem a.$size-12);
		border-radius: var(--_notice-border-radius, a.$border-radius-6);
		background-color: var(--_notice-bg, rgba(255, 255, 255, 0.025));
		align-items: center;

		&__icon {
			display: flex;
			align-items: center;
			justify-content: center;
			margin: 0;
			flex-shrink: 0;
		}

		&__title {
			margin: 0;
			font-size: var(--_notice-font-size, 1rem);
			font-weight: var(--font-weight-medium);
			color: var(--_notice-title-colour, rgba(255, 255, 255, 1));
			line-height: 1;
		}

		&__dismiss {
			cursor: pointer;
			margin-left: auto;
			display: flex;
			align-items: center;
			justify-content: center;
			background: none;
			border: none;
			padding: 0;
			color: var(--_notice-icon-colour);
		}

		:global(svg) {
			min-width: 1.4rem;
			min-height: 1.4rem;
		}
	}
</style>
