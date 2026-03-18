import { CircleCheck, CircleX, Info, TriangleAlert } from '@lucide/svelte';
import type { Component } from 'svelte';
import { writable } from 'svelte/store';

/**
 * Notification type variants.
 */
export type NotificationType = 'info' | 'warning' | 'success' | 'error';

/**
 * Icon configuration for a notification type.
 */
export type IconDetail = {
	icon: Component;
	colour: string;
};

const defaultAlertIcons: Record<NotificationType, IconDetail> = {
	info: { icon: Info, colour: 'var(--colour-semantic-info)' },
	success: { icon: CircleCheck, colour: 'var(--colour-semantic-success)' },
	warning: { icon: TriangleAlert, colour: 'var(--colour-semantic-warning)' },
	error: { icon: CircleX, colour: 'var(--colour-semantic-error)' },
};

/**
 * Writable store holding the active icon map for all Alert and Notice components.
 * Consumers can override individual or all entries via setAlertIcons().
 */
export const alertIconStore = writable<Record<NotificationType, IconDetail>>(defaultAlertIcons);

/**
 * Globally overrides icons for Alert and Notice components.
 * Call this once in your root layout (e.g. +layout.svelte) to supply
 * your own SVG components instead of the default Lucide icons.
 *
 * @example
 * ```ts
 * import { setAlertIcons } from '@ainsleydev/sveltekit-helper';
 * import InfoIcon from '$lib/icons/Info.svelte';
 * import SuccessIcon from '$lib/icons/Success.svelte';
 *
 * setAlertIcons({
 *   info:    { icon: InfoIcon,    colour: 'var(--colour-info)' },
 *   success: { icon: SuccessIcon, colour: 'var(--colour-success)' },
 * });
 * ```
 *
 * @param overrides - Partial map of notification types to icon configurations.
 */
export function setAlertIcons(overrides: Partial<Record<NotificationType, IconDetail>>): void {
	alertIconStore.update((current) => ({ ...current, ...overrides }));
}
