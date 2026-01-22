import { CircleX, Info, CircleCheck, TriangleAlert } from '@lucide/svelte'

import type { AlertType, IconDetail } from './alert.types'

export const alertIcons: Record<AlertType, IconDetail> = {
	info: { icon: Info, colour: 'var(--colour-semantic-info)' },
	success: { icon: CircleCheck, colour: 'var(--colour-semantic-success)' },
	warning: { icon: TriangleAlert, colour: 'var(--colour-semantic-warning)' },
	error: { icon: CircleX, colour: 'var(--colour-semantic-error)' },
}
