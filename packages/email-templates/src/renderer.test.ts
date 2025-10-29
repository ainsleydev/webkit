import * as React from 'react';
import { describe, expect, test } from 'vitest';
import { renderEmail } from './renderer.js';
import { BaseEmail } from './templates/Base.js';
import { generateStyles } from './theme/styles.js';
import type { EmailTheme } from './theme/types.js';

// Example email template for testing.
interface ForgotPasswordProps {
	theme: EmailTheme;
	user: { firstName: string };
	resetUrl: string;
}

const ForgotPasswordEmail = ({ theme, user, resetUrl }: ForgotPasswordProps) => {
	const styles = generateStyles(theme);
	return React.createElement(
		BaseEmail,
		{ theme, previewText: 'Reset your password' },
		React.createElement('h1', { style: styles.heading }, `Hello, ${user.firstName}!`),
		React.createElement(
			'p',
			{ style: styles.text },
			'We received a request to reset your password, please click the button below.',
		),
		React.createElement('a', { href: resetUrl, style: styles.button }, 'Reset Password'),
	);
};

// Example verification email template for testing.
interface VerifyAccountProps {
	theme: EmailTheme;
	user: { firstName: string };
	verifyUrl: string;
}

const VerifyAccountEmail = ({ theme, user, verifyUrl }: VerifyAccountProps) => {
	const styles = generateStyles(theme);
	return React.createElement(
		BaseEmail,
		{ theme, previewText: 'Verify your email' },
		React.createElement('h1', { style: styles.heading }, `Welcome, ${user.firstName}!`),
		React.createElement(
			'p',
			{ style: styles.text },
			'Please verify your email by clicking the button below.',
		),
		React.createElement('a', { href: verifyUrl, style: styles.button }, 'Verify Account'),
	);
};

describe('renderEmail', () => {
	test('renders custom email template', async () => {
		const html = await renderEmail({
			component: ForgotPasswordEmail,
			props: {
				user: { firstName: 'John' },
				resetUrl: 'https://example.com/reset/token123',
			},
		});

		expect(html).toContain('Hello,');
		expect(html).toContain('John');
		expect(html).toContain('reset your password');
		expect(html).toContain('https://example.com/reset/token123');
		expect(html).toContain('Reset Password');
	});

	test('renders different template types', async () => {
		const html = await renderEmail({
			component: VerifyAccountEmail,
			props: {
				user: { firstName: 'Jane' },
				verifyUrl: 'https://example.com/verify/abc123',
			},
		});

		expect(html).toContain('Welcome,');
		expect(html).toContain('Jane');
		expect(html).toContain('verify your email');
		expect(html).toContain('https://example.com/verify/abc123');
		expect(html).toContain('Verify Account');
	});

	test('applies custom theme overrides', async () => {
		const html = await renderEmail({
			component: ForgotPasswordEmail,
			props: {
				user: { firstName: 'Test' },
				resetUrl: 'https://example.com/reset',
			},
			theme: {
				branding: {
					companyName: 'Custom Company',
					logoUrl: 'https://custom.com/logo.png',
					logoWidth: 200,
				},
			},
		});

		expect(html).toContain('Custom Company');
		expect(html).toContain('https://custom.com/logo.png');
		expect(html).toContain('width="200"');
	});

	test('renders plain text when specified', async () => {
		const text = await renderEmail({
			component: ForgotPasswordEmail,
			props: {
				user: { firstName: 'Plain' },
				resetUrl: 'https://example.com/reset',
			},
			plainText: true,
		});

		expect(text).toContain('HELLO, PLAIN!');
		expect(text).toContain('reset your password');
		expect(text).not.toContain('<html>');
		expect(text).not.toContain('<body>');
	});

	test('includes branding footer text', async () => {
		const html = await renderEmail({
			component: VerifyAccountEmail,
			props: {
				user: { firstName: 'Footer' },
				verifyUrl: 'https://example.com/verify',
			},
			theme: {
				branding: {
					footerText: 'Custom footer text.',
				},
			},
		});

		expect(html).toContain('Custom footer text.');
	});

	test('includes website URL in footer when provided', async () => {
		const html = await renderEmail({
			component: ForgotPasswordEmail,
			props: {
				user: { firstName: 'Web' },
				resetUrl: 'https://example.com/reset',
			},
			theme: {
				branding: {
					websiteUrl: 'https://mywebsite.com',
				},
			},
		});

		expect(html).toContain('https://mywebsite.com');
		expect(html).toContain('mywebsite.com');
	});
});
