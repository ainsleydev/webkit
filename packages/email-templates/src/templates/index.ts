import type * as React from 'react';
import { ForgotPasswordEmail } from './ForgotPassword.js';
import type { ForgotPasswordProps } from './ForgotPassword.js';
import { VerifyAccountEmail } from './VerifyAccount.js';
import type { VerifyAccountProps } from './VerifyAccount.js';

/**
 * Template registry mapping template names to their components.
 */
export const templates = {
	'forgot-password': ForgotPasswordEmail,
	'verify-account': VerifyAccountEmail,
} as const;

/**
 * Available template names.
 */
export type TemplateName = keyof typeof templates;

/**
 * Props type mapping for each template.
 */
export type TemplateProps = {
	'forgot-password': ForgotPasswordProps;
	'verify-account': VerifyAccountProps;
};

/**
 * Gets the template component for a given template name.
 *
 * @param name - The template name
 * @returns The React component for the template
 */
export function getTemplate(name: TemplateName): React.FC<TemplateProps[TemplateName]> {
	return templates[name];
}

export { ForgotPasswordEmail, VerifyAccountEmail };
export type { ForgotPasswordProps, VerifyAccountProps };
