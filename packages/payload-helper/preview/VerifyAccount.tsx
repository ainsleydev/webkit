import { renderEmail } from '@ainsleydev/email-templates';
import { VerifyAccountEmail } from '../src/email/VerifyAccountEmail.js';

/**
 * Preview template for the Verify Account email.
 *
 * This file demonstrates how to preview the email with your own branding.
 * Copy this pattern to your own project to preview emails with your actual configuration.
 */
export default async function render() {
	return renderEmail({
		component: VerifyAccountEmail,
		props: {
			user: {
				firstName: 'Jane',
				email: 'jane@example.com',
			},
			verifyUrl: 'https://example.com/admin/users/verify/xyz789token',
			content: {
				previewText: 'Verify your email',
				heading: 'Welcome, Jane!',
				bodyText:
					'Please verify your email by clicking the button below. If you did not create an account, you can safely ignore this email.',
				buttonText: 'Verify Email',
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
