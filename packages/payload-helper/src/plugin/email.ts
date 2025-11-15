import { renderEmail } from '@ainsleydev/email-templates';
import type { Config } from 'payload';

import { ForgotPasswordEmail } from '../email/ForgotPasswordEmail.js';
import { VerifyAccountEmail } from '../email/VerifyAccountEmail.js';
import type { EmailConfig } from '../types.js';

/**
 * Injects email templates into all auth-enabled collections in the Payload config.
 *
 * @param config - The Payload configuration object
 * @param emailConfig - The email configuration from plugin options
 * @returns The modified Payload configuration with email templates injected
 */
export const injectEmailTemplates = (config: Config, emailConfig: EmailConfig): Config => {
	// Get the website URL for branding, defaulting to Payload's serverUrl
	const websiteUrl = emailConfig.frontEndUrl || config.serverURL || '';

	// Merge user theme with websiteUrl for branding
	const themeOverride = {
		...emailConfig.theme,
		branding: {
			...emailConfig.theme?.branding,
			websiteUrl,
		},
	};

	// Find all collections with auth enabled
	const collectionsWithAuth = config.collections?.filter((collection) => collection.auth) || [];

	// If no collections with auth, return config unchanged
	if (collectionsWithAuth.length === 0) {
		return config;
	}

	// Inject email templates into each auth-enabled collection
	const updatedCollections = config.collections?.map((collection) => {
		// Skip collections without auth
		if (!collection.auth) {
			return collection;
		}

		// Clone the collection to avoid mutation
		const updatedCollection = { ...collection };

		// Ensure auth is an object (it could be true or an object)
		if (typeof updatedCollection.auth === 'boolean') {
			updatedCollection.auth = {};
		} else {
			updatedCollection.auth = { ...updatedCollection.auth };
		}

		// Inject forgotPassword email template
		const currentForgotPassword = updatedCollection.auth.forgotPassword;
		updatedCollection.auth.forgotPassword = {
			...(typeof currentForgotPassword === 'object' ? currentForgotPassword : {}),
			generateEmailHTML: async (args) => {
				const token = args?.token || '';
				const user = args?.user || {};
				const resetUrl = `${config.serverURL}/admin/reset/${token}`;

				return renderEmail({
					component: ForgotPasswordEmail,
					props: {
						user: {
							firstName: user?.firstName,
							email: user?.email,
						},
						resetUrl,
						content: emailConfig.forgotPassword,
					},
					theme: themeOverride,
				});
			},
		};

		// Inject verify email template
		const currentVerify = updatedCollection.auth.verify;
		updatedCollection.auth.verify = {
			...(typeof currentVerify === 'object' ? currentVerify : {}),
			generateEmailHTML: async (args) => {
				const token = args?.token || '';
				const user = args?.user || {};
				// For verify emails, the token is used in the verification URL
				const verifyUrl = `${config.serverURL}/admin/${collection.slug}/verify/${token}`;

				return renderEmail({
					component: VerifyAccountEmail,
					props: {
						user: {
							firstName: user?.firstName,
							email: user?.email,
						},
						verifyUrl,
						content: emailConfig.verifyAccount,
					},
					theme: themeOverride,
				});
			},
		};

		return updatedCollection;
	});

	return {
		...config,
		collections: updatedCollections,
	};
};
