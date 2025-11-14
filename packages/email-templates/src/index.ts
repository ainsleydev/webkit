/**
 * @ainsleydev/email-templates
 *
 * Composable email template building blocks with theme system for React Email.
 * Create your own email templates using the BaseEmail component and theme system.
 */

// Re-exported React Email components for convenience.
export * from '@react-email/components';
export type { RenderEmailOptions } from './renderer.js';
// Main rendering function.
export { renderEmail } from './renderer.js';
// Base email component for building templates.
export { BaseEmail } from './templates/index.js';
export type { EmailBranding, EmailColours, EmailTheme, PartialEmailTheme } from './theme/index.js';
// Theme system.
export { defaultTheme, mergeTheme } from './theme/index.js';
export { generateStyles } from './theme/styles.js';
