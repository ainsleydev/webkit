import { type Icon as IconType } from '@lucide/svelte'

export type AlertType = 'info' | 'warning' | 'success' | 'error'

export type IconDetail = {
	icon: typeof IconType
	colour: string
}
