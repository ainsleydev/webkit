import type { Config } from 'payload';
import type { AdminLogoConfig } from '../types.js';

/**
 * Injects the admin Logo component into the Payload config.
 */
export const injectAdminLogo = (
	config: Config,
	logoConfig: AdminLogoConfig,
	siteName: string,
): Config => ({
	...config,
	admin: {
		...config.admin,
		components: {
			...config.admin?.components,
			graphics: {
				...config.admin?.components?.graphics,
				Logo: {
					path: '@ainsleydev/payload-helper/dist/admin/components/Logo',
					exportName: 'Logo',
					clientProps: {
						config: {
							...logoConfig,
							alt: logoConfig.alt || siteName,
						},
					},
				},
			},
		},
	},
});
