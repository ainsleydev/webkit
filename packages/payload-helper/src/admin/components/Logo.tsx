'use client';

import { useTheme } from '@payloadcms/ui';
// @ts-expect-error - next/image is a peer dependency and will be available in consumer projects
import Image from 'next/image';
import type React from 'react';

/**
 * Configuration for the Logo component.
 */
export type LogoConfig = {
	path: string;
	darkModePath?: string;
	width?: number;
	height?: number;
	alt: string;
	className?: string;
};

/**
 * Props for the Logo component.
 */
export type LogoProps = {
	config: LogoConfig;
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
				alt={config.alt}
				width={config.width || 150}
				height={config.height || 40}
			/>
		</figure>
	);
};
