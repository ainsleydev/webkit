import { Heading, Text, Button, Section } from '@react-email/components';
import * as React from 'react';
import { BaseEmail } from './Base.js';
import type { EmailTheme } from '../theme/types.js';
import { generateStyles } from '../theme/styles.js';

export interface VerifyAccountProps {
	theme: EmailTheme;
	user: { firstName: string };
	verifyUrl: string;
}

/**
 * Account verification email template.
 * Sends a verification link to new users.
 */
export const VerifyAccountEmail = ({ theme, user, verifyUrl }: VerifyAccountProps) => {
	const styles = generateStyles(theme);

	return (
		<BaseEmail previewText='Verify your email' theme={theme}>
			<Heading style={styles.heading}>Welcome, {user.firstName}!</Heading>
			<Text style={styles.text}>
				Please verify your email by clicking the button below. If you did not create an
				account, you can safely ignore this email.
			</Text>
			<Section style={{ textAlign: 'center' }}>
				<Button href={verifyUrl} style={styles.button}>
					Verify Account
				</Button>
			</Section>
		</BaseEmail>
	);
};
