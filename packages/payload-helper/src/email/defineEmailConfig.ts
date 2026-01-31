import type { EmailConfig } from '../types.js';

/**
 * Helper function for defining email configuration with type safety.
 * Returns the config as-is — this is an identity function that provides
 * autocomplete and type checking in config files.
 *
 * @example
 * ```typescript
 * // src/email.config.ts
 * import { defineEmailConfig } from '@ainsleydev/payload-helper';
 *
 * export default defineEmailConfig({
 *   frontEndUrl: 'https://yoursite.com',
 *   theme: {
 *     branding: {
 *       companyName: 'My Company',
 *       logoUrl: 'https://yoursite.com/logo.png',
 *     },
 *   },
 * });
 * ```
 */
export const defineEmailConfig = (config: EmailConfig): EmailConfig => config;
