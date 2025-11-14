import { renderEmail } from '@ainsleydev/email-templates';
import { ForgotPasswordEmail } from '../src/email/ForgotPasswordEmail.js';

/**
 * Preview template for the Forgot Password email.
 *
 * This file demonstrates how to preview the email with your own branding.
 * Copy this pattern to your own project to preview emails with your actual configuration.
 */
export default async function render() {
	return renderEmail({
		component: ForgotPasswordEmail,
		props: {
			user: {
				firstName: 'John',
				email: 'john@example.com',
			},
			resetUrl: 'https://example.com/admin/reset/abc123token',
			content: {
				previewText: 'Reset your password',
				heading: 'Hello, John!',
				bodyText:
					'We received a request to reset your password, please click the button below. If you did not request a password reset, you can safely ignore this email.',
				buttonText: 'Reset Password',
			},
		},
		theme: {
			branding: {
				companyName: 'My Company',
				logoUrl: 'https://via.placeholder.com/150x40/ff5043/ffffff?text=My+Company',
				websiteUrl: 'https://example.com',
			},
			colours: {
				background: {
					accent: '#ff5043',
				},
			},
		},
	});
}
