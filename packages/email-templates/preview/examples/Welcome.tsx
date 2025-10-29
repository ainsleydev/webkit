import { Button, Heading, Section, Text } from '@react-email/components';
import * as React from 'react';
import { BaseEmail } from '../../src/templates/Base.js';
import type { EmailTheme } from '../../src/theme/types.js';

interface WelcomeEmailProps {
	theme: EmailTheme;
	userName?: string;
	loginUrl?: string;
}

/**
 * Welcome email template for new user onboarding.
 */
export const WelcomeEmail = ({
	theme,
	userName = 'User',
	loginUrl = 'https://example.com/login',
}: WelcomeEmailProps) => {
	return (
		<BaseEmail theme={theme} previewText={`Welcome to ${theme.branding.companyName}!`}>
			<Heading
				style={{
					color: theme.colours.text.heading,
					fontSize: '24px',
					fontWeight: 'bold',
					marginBottom: '20px',
				}}
			>
				Welcome to {theme.branding.companyName}!
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
				Thank you for registering an account with us. We're excited to have you on board!
			</Text>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '16px',
					lineHeight: '24px',
					marginBottom: '30px',
				}}
			>
				Your account has been successfully created and you can now access all our features.
				Click the button below to get started.
			</Text>

			<Section style={{ textAlign: 'center', marginBottom: '30px' }}>
				<Button
					href={loginUrl}
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
					Get Started
				</Button>
			</Section>

			<Text
				style={{
					color: theme.colours.text.body,
					fontSize: '14px',
					lineHeight: '20px',
					marginTop: '30px',
				}}
			>
				If you have any questions, feel free to reach out to our support team.
			</Text>
		</BaseEmail>
	);
};
