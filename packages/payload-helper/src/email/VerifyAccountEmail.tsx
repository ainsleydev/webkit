import { BaseEmail, Button, Heading, Section, Text } from '@ainsleydev/email-templates';
import type { EmailTheme } from '@ainsleydev/email-templates';
import * as React from 'react';

/**
 * Props for the VerifyAccountEmail component.
 */
export interface VerifyAccountEmailProps {
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
	 * The URL for verifying the account.
	 */
	verifyUrl: string;

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
 * Email template for account verification in Payload CMS.
 *
 * @param props - The component props
 * @returns The rendered email component
 */
export const VerifyAccountEmail = ({
	theme,
	user,
	verifyUrl,
	content,
}: VerifyAccountEmailProps) => {
	const userName = user.firstName || user.email || 'there';
	const previewText = content?.previewText || 'Verify your email';
	const heading = content?.heading || `Welcome, ${userName}!`;
	const bodyText =
		content?.bodyText ||
		'Please verify your email by clicking the button below. If you did not request a password reset, you can safely ignore this email.';
	const buttonText = content?.buttonText || 'Verify Email';

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
					href={verifyUrl}
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
