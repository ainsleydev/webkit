<script lang="ts" module>
	import { type Icon as IconType } from '@lucide/svelte'
	import type { Snippet } from 'svelte'

	export type AlertType = 'info' | 'warning' | 'success' | 'error'

	export type AlertProps = {
		type?: AlertType
		title?: string
		children?: Snippet
		visible?: boolean
		dismiss?: boolean
		icon?: typeof IconType
	}
</script>

<script lang="ts">
	import { X } from '@lucide/svelte'
	import { fade } from 'svelte/transition'

	import { alertIcons } from './alertIcons'

	let {
		type = 'info',
		title = '',
		children,
		visible = true,
		dismiss = false,
		icon: customIcon,
		...restProps
	}: AlertProps = $props()

	const iconDetail = $derived(alertIcons[type])
	const Icon = $derived(customIcon || iconDetail.icon)
	const hide = () => (visible = false)
</script>

<!--
	@component

	Full-width alert component for displaying important messages with optional body text.
	Supports title, children content, and custom icons with dismissible functionality.

	@example
	```svelte
	<Alert type="info" title="New features available">
		Check out the latest updates in your dashboard.
	</Alert>

	<Alert type="warning" title="Maintenance scheduled" dismiss />

	<Alert type="error" title="Payment failed" dismiss>
		Your card was declined. Please update your payment method.
	</Alert>
	```
-->
{#if visible}
	<div
		class="alert alert--{type}"
		role="alert"
		style="--_alert-icon-colour: {iconDetail.colour}"
		transition:fade={{ duration: 300 }}
		{...restProps}
	>
		<!-- Icon -->
		<figure class="alert__icon">
			<Icon size={24} color={iconDetail.colour} strokeWidth={1.2}></Icon>
		</figure>
		<!-- Content -->
		<div class="alert__content">
			{#if title}
				<p class="alert__title">
					{title}
				</p>
			{/if}
			{#if children}
				<p class="alert__text">
					{@render children()}
				</p>
			{/if}
		</div>
		<!-- Dismiss -->
		{#if dismiss}
			<button class="alert__dismiss" onclick={hide} aria-label="Close">
				<X size={24} color={iconDetail.colour} />
			</button>
		{/if}
	</div>
{/if}

<style lang="scss">
	@use '$lib/scss/abstracts' as a;

	.alert {
		$self: &;

		position: relative;
		display: flex;
		gap: var(--_alert-gap, a.$size-12);
		width: 100%;
		border-radius: var(--_alert-border-radius, a.$border-radius-6);
		padding: var(--_alert-padding, a.$size-24);
		background-color: var(--_alert-bg, rgba(255, 255, 255, 0.02));

		&__icon {
			display: flex;
			align-items: center;
			justify-content: center;
			margin: 0;
			flex-shrink: 0;
		}

		&__content {
			display: grid;
			gap: var(--_alert-content-gap, a.$size-8);
		}

		&__title {
			font-weight: var(--_alert-title-font-weight, var(--font-weight-semibold));
			margin: 0;
			line-height: 1;
			color: var(--_alert-title-colour, rgba(255, 255, 255, 1));

			&:empty {
				display: none;
			}
		}

		&__text {
			margin: 0;
			line-height: var(--_alert-text-line-height, 1.4);
			color: var(--_alert-text-colour, rgba(255, 255, 255, 50%));
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
			color: var(--_alert-icon-colour);
		}

		:global(svg) {
			min-width: 1.4rem;
			min-height: 1.4rem;
		}
	}
</style>
