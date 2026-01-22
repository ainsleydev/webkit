import { CircleX, Info, CircleCheck, TriangleAlert, type Icon as IconType } from '@lucide/svelte'

/**
 * Notification type variants
 */
type NotificationType = 'info' | 'warning' | 'success' | 'error'

/**
 * Icon configuration for notification types
 */
type IconDetail = {
	icon: typeof IconType
	colour: string
}

export const alertIcons: Record<NotificationType, IconDetail> = {
	info: { icon: Info, colour: 'var(--colour-semantic-info)' },
	success: { icon: CircleCheck, colour: 'var(--colour-semantic-success)' },
	warning: { icon: TriangleAlert, colour: 'var(--colour-semantic-warning)' },
	error: { icon: CircleX, colour: 'var(--colour-semantic-error)' },
}
