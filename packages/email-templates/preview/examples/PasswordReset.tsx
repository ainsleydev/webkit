import { Button, Heading, Section, Text } from '@react-email/components';
// biome-ignore lint/style/useImportType: React is needed for JSX
import * as React from 'react';
import { BaseEmail } from '../../src/templates/Base.js';
import type { EmailTheme } from '../../src/theme/types.js';

interface PasswordResetEmailProps {
	theme: EmailTheme;
	userName?: string;
	resetUrl?: string;
	expiryMinutes?: number;
}

/**
 * Password reset email template for authentication flow.
 */
export const PasswordResetEmail = ({
	theme,
	userName = 'User',
	resetUrl = 'https://example.com/reset-password?token=abc123',
	expiryMinutes = 60,
}: PasswordResetEmailProps) => {
	return (
		<BaseEmail theme={theme} previewText='Reset your password'>
			<Heading
				style={{
					color: theme.colours.text.heading,
					fontSize: '24px',
					fontWeight: 'bold',
					marginBottom: '20px',
				}}
			>
				Password Reset Request
			</Heading>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '16px',
					lineHeight: '24px',
					marginBottom: '20px',
				}}
			>
				Hi {userName},
			</Text>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '16px',
					lineHeight: '24px',
					marginBottom: '20px',
				}}
			>
				We received a request to reset the password for your account. If you didn't make
				this request, you can safely ignore this email.
			</Text>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '16px',
					lineHeight: '24px',
					marginBottom: '30px',
				}}
			>
				To reset your password, click the button below. This link will expire in{' '}
				{expiryMinutes} minutes.
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
					Reset Password
				</Button>
			</Section>

			<Section
				style={{
					backgroundColor: theme.colours.background.highlight,
					padding: '15px',
					borderRadius: '5px',
					borderLeft: `4px solid ${theme.colours.text.negative}`,
					marginBottom: '20px',
				}}
			>
				<Text
					style={{
						color: theme.colours.text.body,
						fontSize: '14px',
						lineHeight: '20px',
						margin: '0',
					}}
				>
					<strong>Security tip:</strong> If you didn't request a password reset, please
					contact our support team immediately.
				</Text>
			</Section>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '14px',
					lineHeight: '20px',
					marginTop: '30px',
				}}
			>
				If the button doesn't work, you can copy and paste this link into your browser:
			</Text>

			<Text
				style={{
					color: theme.colours.text.action,
					fontSize: '12px',
					lineHeight: '18px',
					wordBreak: 'break-all',
				}}
			>
				{resetUrl}
			</Text>
		</BaseEmail>
	);
};
