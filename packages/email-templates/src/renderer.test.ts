import { describe, test, expect } from 'vitest';
import { renderEmail } from './renderer.js';

describe('renderEmail', () => {
	test('renders forgot-password template', async () => {
		const html = await renderEmail({
			template: 'forgot-password',
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

	test('renders verify-account template', async () => {
		const html = await renderEmail({
			template: 'verify-account',
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
			template: 'forgot-password',
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
			template: 'forgot-password',
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
			template: 'verify-account',
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
			template: 'forgot-password',
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
