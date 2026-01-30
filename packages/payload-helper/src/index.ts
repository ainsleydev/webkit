// Main plugin
export { payloadHelper } from './plugin.js';

// Types
export type {
	PayloadHelperPluginConfig,
	AdminConfig,
	AdminIconConfig,
	AdminLogoConfig,
	EmailConfig,
	EmailContentOverrides,
	SettingsConfig,
	WebServerConfig,
} from './types.js';

// Collections
export type { MediaArgs } from './collections/Media.js';
export { Media, imageSizes, imageSizesWithAvif, Redirects } from './collections/index.js';

// Globals
export type { SettingsArgs } from './globals/Settings.js';
export { Settings, Navigation, countries, languages } from './globals/index.js';

// Email Components
export { ForgotPasswordEmail } from './email/ForgotPasswordEmail.js';
export type { ForgotPasswordEmailProps } from './email/ForgotPasswordEmail.js';
export { VerifyAccountEmail } from './email/VerifyAccountEmail.js';
export type { VerifyAccountEmailProps } from './email/VerifyAccountEmail.js';

// Admin Components
export type { IconProps } from './admin/components/Icon.js';
export type { LogoProps } from './admin/components/Logo.js';

// Utilities
export {
	env,
	fieldHasName,
	validateURL,
	validatePostcode,
	htmlToLexical,
	lexicalToHtml,
} from './util/index.js';

// Common/Reusable
export { SEOFields } from './common/index.js';

// Email Config Helper
export { defineEmailConfig } from './email/defineEmailConfig.js';

// Endpoints
export { findBySlug } from './endpoints/index.js';

// Schema utilities
export type { SchemaOptions } from './plugin/schema.js';
export { fieldMapper, schemas, addGoJSONSchema } from './plugin/schema.js';
