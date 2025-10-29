import { render } from '@react-email/render';
import * as React from 'react';
import type { PartialEmailTheme } from './theme/types.js';
import { mergeTheme } from './theme/merge.js';
import { getTemplate, type TemplateName, type TemplateProps } from './templates/index.js';

/**
 * Options for rendering an email template.
 */
export interface RenderEmailOptions<T extends TemplateName> {
	/**
	 * The name of the template to render.
	 */
	template: T;
	/**
	 * Props specific to the template being rendered.
	 */
	props: Omit<TemplateProps[T], 'theme'>;
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
 * Renders an email template to HTML string.
 *
 * @param options - Rendering options including template name, props, and theme
 * @returns HTML string ready to be sent via email service
 *
 * @example
 * ```typescript
 * const html = await renderEmail({
 *   template: 'forgot-password',
 *   props: {
 *     user: { firstName: 'John' },
 *     resetUrl: 'https://example.com/reset/token123'
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
export async function renderEmail<T extends TemplateName>(
	options: RenderEmailOptions<T>,
): Promise<string> {
	const { template, props, theme: partialTheme, plainText = false } = options;

	// Merge partial theme with defaults.
	const theme = mergeTheme(partialTheme);

	// Get the template component.
	const TemplateComponent = getTemplate(template);

	// Create the React element with merged theme.
	const element = React.createElement(TemplateComponent, {
		...props,
		theme,
	});

	// Render to HTML or plain text.
	return render(element, { plainText });
}
