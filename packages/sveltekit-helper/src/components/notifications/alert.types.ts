import { type Icon as IconType } from '@lucide/svelte'

/** Available notification type variants */
export type AlertType = 'info' | 'warning' | 'success' | 'error'

/** Icon configuration for notification types */
export type IconDetail = {
	icon: typeof IconType
	colour: string
}
