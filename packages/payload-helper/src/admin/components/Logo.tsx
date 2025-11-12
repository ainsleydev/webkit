'use client';

import { useTheme } from '@payloadcms/ui';
// @ts-expect-error - next/image is a peer dependency and will be available in consumer projects
import Image from 'next/image';
import type React from 'react';
import type { AdminLogoConfig } from '../../types.js';

/**
 * Props for the Logo component.
 */
export type LogoProps = {
	config: AdminLogoConfig;
};

/**
 * Logo component that displays a configurable logo with theme support.
 *
 * @param {LogoProps} props - Component props containing logo configuration
 * @returns {React.ReactElement} Logo component
 */
export const Logo: React.FC<LogoProps> = ({ config }) => {
	const { theme } = useTheme();
	const imagePath = theme === 'light' || !config.darkModePath ? config.path : config.darkModePath;

	return (
		<figure className={config.className ? config.className : undefined}>
			<Image
				src={imagePath}
				alt={config.alt || 'Logo'}
				width={config.width || 150}
				height={config.height || 40}
			/>
		</figure>
	);
};
