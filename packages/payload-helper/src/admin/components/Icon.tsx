'use client';

import { useTheme } from '@payloadcms/ui';
// @ts-expect-error - next/image is a peer dependency and will be available in consumer projects
import Image from 'next/image';
import type React from 'react';
import type { AdminIconConfig } from '../../types.js';

/**
 * Props for the Icon component.
 */
export type IconProps = {
	config: AdminIconConfig;
};

/**
 * Icon component that displays a configurable icon with theme support.
 * The icon appears in the top left corner of the Payload dashboard.
 *
 * @param {IconProps} props - Component props containing icon configuration
 * @returns {React.ReactElement} Icon component
 */
export const Icon: React.FC<IconProps> = ({ config }) => {
	const { theme } = useTheme();
	const imagePath = theme === 'light' || !config.darkModePath ? config.path : config.darkModePath;

	return (
		<figure className={config.className ? config.className : undefined}>
			<Image
				src={imagePath}
				alt={config.alt || 'Icon'}
				width={config.width || 120}
				height={config.height || 120}
			/>
		</figure>
	);
};
