import type { Config } from 'payload';
import type { AdminIconConfig } from '../types.js';

/**
 * Injects the admin Icon component into the Payload config.
 */
export const injectAdminIcon = (
	config: Config,
	iconConfig: AdminIconConfig,
	siteName: string,
): Config => ({
	...config,
	admin: {
		...config.admin,
		components: {
			...config.admin?.components,
			graphics: {
				...config.admin?.components?.graphics,
				Icon: {
					path: '@ainsleydev/payload-helper/dist/admin/components/Icon',
					exportName: 'Icon',
					clientProps: {
						config: {
							...iconConfig,
							alt: iconConfig.alt || siteName,
						},
					},
				},
			},
		},
	},
});
