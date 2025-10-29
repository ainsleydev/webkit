import { render } from '@react-email/render';
import * as React from 'react';
import type { EmailTheme, PartialEmailTheme } from './theme/types.js';
import { mergeTheme } from './theme/merge.js';

/**
 * Options for rendering an email template.
 */
export interface RenderEmailOptions<P = Record<string, unknown>> {
	/**
	 * The React component to render.
	 */
	component: React.ComponentType<P & { theme: EmailTheme }>;
	/**
	 * Props to pass to the component (excluding theme).
	 */
	props: Omit<P, 'theme'>;
	/**
	 * Optional theme overrides. Will be merged with default theme.
	 */
	theme?: PartialEmailTheme;
	/**
	 * Whether to render as plain text instead of HTML.
	 * @default false
	 */
	plainText?: boolean;
}

/**
 * Renders an email template component to HTML string.
 *
 * @param options - Rendering options including component, props, and theme
 * @returns HTML string ready to be sent via email service
 *
 * @example
 * ```typescript
 * import { renderEmail } from '@ainsleydev/email-templates'
 * import { MyEmailTemplate } from './emails/MyTemplate'
 *
 * const html = await renderEmail({
 *   component: MyEmailTemplate,
 *   props: {
 *     user: { firstName: 'John' },
 *     actionUrl: 'https://example.com/action'
 *   },
 *   theme: {
 *     branding: {
 *       companyName: 'My Company',
 *       logoUrl: 'https://example.com/logo.png'
 *     }
 *   }
 * })
 * ```
 */
export async function renderEmail<P = Record<string, unknown>>(
	options: RenderEmailOptions<P>,
): Promise<string> {
	const { component: Component, props, theme: partialTheme, plainText = false } = options;

	// Merge partial theme with defaults.
	const theme = mergeTheme(partialTheme);

	// Create the React element with merged theme.
	const element = React.createElement(Component, {
		...props,
		theme,
	} as P & { theme: EmailTheme });

	// Render to HTML or plain text.
	return render(element, { plainText });
}
