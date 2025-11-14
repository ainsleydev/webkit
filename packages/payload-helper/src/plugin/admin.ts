import type { Config } from 'payload';
import type { AdminIconConfig, AdminLogoConfig } from '../types.js';

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
