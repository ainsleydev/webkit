import { BaseEmail, Button, Heading, Section, Text } from '@ainsleydev/email-templates';
import type { EmailTheme } from '@ainsleydev/email-templates';
import * as React from 'react';

/**
 * Props for the ForgotPasswordEmail component.
 */
export interface ForgotPasswordEmailProps {
	/**
	 * The email theme (required by renderEmail).
	 */
	theme: EmailTheme;

	/**
	 * The user object containing user information.
	 */
	user: {
		firstName?: string;
		email?: string;
	};

	/**
	 * The URL for resetting the password.
	 */
	resetUrl: string;

	/**
	 * Optional content overrides.
	 */
	content?: {
		previewText?: string;
		heading?: string;
		bodyText?: string;
		buttonText?: string;
	};
}

/**
 * Email template for password reset requests in Payload CMS.
 *
 * @param props - The component props
 * @returns The rendered email component
 */
export const ForgotPasswordEmail = ({
	theme,
	user,
	resetUrl,
	content,
}: ForgotPasswordEmailProps) => {
	const userName = user.firstName || user.email || 'there';
	const previewText = content?.previewText || 'Reset your password';
	const heading = content?.heading || `Hello, ${userName}!`;
	const bodyText =
		content?.bodyText ||
		'We received a request to reset your password, please click the button below. If you did not request a password reset, you can safely ignore this email.';
	const buttonText = content?.buttonText || 'Reset Password';

	return (
		<BaseEmail theme={theme} previewText={previewText}>
			<Heading
				style={{
					color: theme.colours.text.heading,
					fontSize: '24px',
					fontWeight: 'bold',
					marginBottom: '20px',
				}}
			>
				{heading}
			</Heading>
			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '16px',
					lineHeight: '24px',
					marginBottom: '30px',
				}}
			>
				{bodyText}
			</Text>
			<Section style={{ textAlign: 'center', marginBottom: '30px' }}>
				<Button
					href={resetUrl}
					style={{
						backgroundColor: theme.colours.background.accent,
						color: theme.colours.text.heading,
						padding: '12px 32px',
						borderRadius: '5px',
						fontSize: '16px',
						fontWeight: 'bold',
						textDecoration: 'none',
						display: 'inline-block',
					}}
				>
					{buttonText}
				</Button>
			</Section>
		</BaseEmail>
	);
};
