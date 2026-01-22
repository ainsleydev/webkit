import { type Icon as IconType } from '@lucide/svelte'

/**
 * Available alert type variants
 * - info: Informational message
 * - warning: Warning or cautionary message
 * - success: Success confirmation message
 * - error: Error or critical message
 */
export type AlertType = 'info' | 'warning' | 'success' | 'error'

/**
 * Icon configuration for alert types
 * Contains the Lucide icon component and corresponding colour
 */
export type IconDetail = {
	/** Lucide icon component to display */
	icon: typeof IconType
	/** CSS colour value for the icon (usually a semantic colour variable) */
	colour: string
}
