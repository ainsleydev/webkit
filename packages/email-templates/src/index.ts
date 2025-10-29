/**
 * @ainsleydev/email-templates
 *
 * Composable, reusable email templates built with React Email.
 * Works with JavaScript (primary) and Go (via CLI).
 */

// Main rendering function.
export { renderEmail } from './renderer.js';
export type { RenderEmailOptions } from './renderer.js';

// Theme system.
export { defaultTheme, mergeTheme } from './theme/index.js';
export type { EmailTheme, EmailColours, EmailBranding, PartialEmailTheme } from './theme/index.js';

// Templates.
export { ForgotPasswordEmail, VerifyAccountEmail } from './templates/index.js';
export type {
	TemplateName,
	TemplateProps,
	ForgotPasswordProps,
	VerifyAccountProps,
} from './templates/index.js';
