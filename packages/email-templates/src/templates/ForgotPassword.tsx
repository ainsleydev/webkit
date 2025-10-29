import { Heading, Text, Button, Section } from '@react-email/components';
import * as React from 'react';
import { BaseEmail } from './Base.js';
import type { EmailTheme } from '../theme/types.js';
import { generateStyles } from '../theme/styles.js';

export interface ForgotPasswordProps {
	theme: EmailTheme;
	user: { firstName: string };
	resetUrl: string;
}

/**
 * Forgot password email template.
 * Sends a password reset link to the user.
 */
export const ForgotPasswordEmail = ({ theme, user, resetUrl }: ForgotPasswordProps) => {
	const styles = generateStyles(theme);

	return (
		<BaseEmail previewText='Reset your password' theme={theme}>
			<Heading style={styles.heading}>Hello, {user.firstName}!</Heading>
			<Text style={styles.text}>
				We received a request to reset your password, please click the button below. If you
				did not request a password reset, you can safely ignore this email.
			</Text>
			<Section style={{ textAlign: 'center' }}>
				<Button href={resetUrl} style={styles.button}>
					Reset Password
				</Button>
			</Section>
		</BaseEmail>
	);
};
