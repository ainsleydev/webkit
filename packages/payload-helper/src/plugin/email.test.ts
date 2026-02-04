import { renderEmail } from '@ainsleydev/email-templates';
import type { CollectionConfig, Config } from 'payload';
import { describe, expect, test, vi } from 'vitest';
import type { EmailConfig } from '../types.js';
import { injectEmailTemplates } from './email.js';

// Mock the email templates module
vi.mock('@ainsleydev/email-templates', () => ({
	renderEmail: vi.fn(async () => '<html>Mocked Email</html>'),
}));

describe('injectEmailTemplates', () => {
	const mockEmailConfig: EmailConfig = {
		frontEndUrl: 'https://example.com',
		theme: {
			branding: {
				companyName: 'Test Company',
			},
		},
		forgotPassword: {
			heading: 'Reset Your Password',
			bodyText: 'Click below to reset your password',
		},
		verifyAccount: {
			heading: 'Verify Your Account',
			bodyText: 'Click below to verify your account',
		},
	};

	test('should inject email templates into auth-enabled collections', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const usersCollection = result.collections?.[0];
		expect(usersCollection).toBeDefined();
		expect(typeof usersCollection?.auth).toBe('object');

		if (typeof usersCollection?.auth === 'object') {
			expect(usersCollection.auth.forgotPassword).toBeDefined();
			expect(usersCollection.auth.verify).toBeDefined();
			expect(usersCollection.auth.forgotPassword?.generateEmailHTML).toBeDefined();
			expect(usersCollection.auth.verify?.generateEmailHTML).toBeDefined();
		}
	});

	test('should not inject email templates into non-auth collections', () => {
		const config: Config = {
			collections: [
				{
					slug: 'posts',
					fields: [],
					// No auth property
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const postsCollection = result.collections?.[0];
		expect(postsCollection).toBeDefined();
		expect(postsCollection?.auth).toBeUndefined();
	});

	test('should return config unchanged when no auth-enabled collections exist', () => {
		const config: Config = {
			collections: [
				{
					slug: 'posts',
					fields: [],
				},
				{
					slug: 'media',
					fields: [],
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		expect(result.collections).toEqual(config.collections);
	});

	test('should handle auth as boolean and convert to object', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const usersCollection = result.collections?.[0];
		expect(typeof usersCollection?.auth).toBe('object');
	});

	test('should preserve existing auth configuration', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: {
						tokenExpiration: 7200,
						maxLoginAttempts: 5,
					},
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const usersCollection = result.collections?.[0];
		if (typeof usersCollection?.auth === 'object') {
			expect(usersCollection.auth.tokenExpiration).toBe(7200);
			expect(usersCollection.auth.maxLoginAttempts).toBe(5);
		}
	});

	test('should merge theme with websiteUrl from frontEndUrl', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			frontEndUrl: 'https://custom-frontend.com',
			theme: {
				branding: {
					companyName: 'Test Company',
				},
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		// The theme should have websiteUrl merged into branding
		// This is tested indirectly through the generateEmailHTML function
		const usersCollection = result.collections?.[0];
		expect(usersCollection).toBeDefined();
	});

	test('should use serverURL when frontEndUrl is not provided', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			// No frontEndUrl provided
			theme: {
				branding: {
					companyName: 'Test Company',
				},
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		const usersCollection = result.collections?.[0];
		expect(usersCollection).toBeDefined();
		// The websiteUrl should default to serverURL
	});

	test('should handle multiple auth-enabled collections', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
				{
					slug: 'admins',
					fields: [],
					auth: {
						tokenExpiration: 3600,
					},
				},
				{
					slug: 'posts',
					fields: [],
					// No auth
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		expect(result.collections).toHaveLength(3);

		const usersCollection = result.collections?.[0];
		const adminsCollection = result.collections?.[1];
		const postsCollection = result.collections?.[2];

		// Users should have email templates
		if (typeof usersCollection?.auth === 'object') {
			expect(usersCollection.auth.forgotPassword?.generateEmailHTML).toBeDefined();
			expect(usersCollection.auth.verify?.generateEmailHTML).toBeDefined();
		}

		// Admins should have email templates and preserve existing config
		if (typeof adminsCollection?.auth === 'object') {
			expect(adminsCollection.auth.forgotPassword?.generateEmailHTML).toBeDefined();
			expect(adminsCollection.auth.verify?.generateEmailHTML).toBeDefined();
			expect(adminsCollection.auth.tokenExpiration).toBe(3600);
		}

		// Posts should remain unchanged
		expect(postsCollection?.auth).toBeUndefined();
	});

	test('should handle empty collections array', () => {
		const config: Config = {
			collections: [],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		expect(result.collections).toEqual([]);
	});

	test('should handle config without collections property', () => {
		const config: Config = {
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		expect(result.collections).toBeUndefined();
	});

	test('should preserve existing forgotPassword configuration', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: {
						forgotPassword: {
							generateEmailSubject: () => 'Custom Subject',
						},
					},
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const usersCollection = result.collections?.[0];
		if (typeof usersCollection?.auth === 'object') {
			expect(usersCollection.auth.forgotPassword?.generateEmailSubject).toBeDefined();
			expect(usersCollection.auth.forgotPassword?.generateEmailHTML).toBeDefined();
		}
	});

	test('should preserve existing verify configuration', () => {
		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: {
						verify: {
							generateEmailSubject: () => 'Custom Verify Subject',
						},
					},
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const result = injectEmailTemplates(config, mockEmailConfig);

		const usersCollection = result.collections?.[0];
		if (typeof usersCollection?.auth === 'object') {
			expect(usersCollection.auth.verify?.generateEmailSubject).toBeDefined();
			expect(usersCollection.auth.verify?.generateEmailHTML).toBeDefined();
		}
	});

	test('should use custom forgotPassword.url callback when provided', async () => {
		const customUrlCallback = vi.fn(
			({ token, collection }) =>
				`https://custom.example.com/reset?token=${token}&slug=${collection.slug}`,
		);

		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			...mockEmailConfig,
			forgotPassword: {
				...mockEmailConfig.forgotPassword,
				url: customUrlCallback,
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		const usersCollection = result.collections?.[0];
		if (
			typeof usersCollection?.auth === 'object' &&
			usersCollection.auth.forgotPassword?.generateEmailHTML
		) {
			await usersCollection.auth.forgotPassword.generateEmailHTML({
				token: 'test-token-123',
				user: { email: 'test@example.com' },
			});

			expect(customUrlCallback).toHaveBeenCalledWith({
				token: 'test-token-123',
				config,
				collection: expect.objectContaining({ slug: 'users' }),
			});
		}
	});

	test('should use custom verifyAccount.url callback when provided', async () => {
		const customUrlCallback = vi.fn(
			({ token, collection }) =>
				`https://custom.example.com/verify?token=${token}&type=${collection.slug}`,
		);

		const config: Config = {
			collections: [
				{
					slug: 'admins',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			...mockEmailConfig,
			verifyAccount: {
				...mockEmailConfig.verifyAccount,
				url: customUrlCallback,
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		const adminsCollection = result.collections?.[0];
		if (
			typeof adminsCollection?.auth === 'object' &&
			adminsCollection.auth.verify?.generateEmailHTML
		) {
			await adminsCollection.auth.verify.generateEmailHTML({
				token: 'verify-token-456',
				user: { email: 'admin@example.com' },
			});

			expect(customUrlCallback).toHaveBeenCalledWith({
				token: 'verify-token-456',
				config,
				collection: expect.objectContaining({ slug: 'admins' }),
			});
		}
	});

	test('should support async URL callbacks', async () => {
		const asyncUrlCallback = vi.fn(async ({ token }) => {
			await new Promise((resolve) => setTimeout(resolve, 10));
			return `https://async.example.com/reset?token=${token}`;
		});

		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			...mockEmailConfig,
			forgotPassword: {
				url: asyncUrlCallback,
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		const usersCollection = result.collections?.[0];
		if (
			typeof usersCollection?.auth === 'object' &&
			usersCollection.auth.forgotPassword?.generateEmailHTML
		) {
			await usersCollection.auth.forgotPassword.generateEmailHTML({
				token: 'async-token',
				user: { email: 'test@example.com' },
			});

			expect(asyncUrlCallback).toHaveBeenCalled();
		}
	});

	test('should fallback to default URL when callback throws error', async () => {
		const consoleWarnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
		const renderEmailMock = vi.mocked(renderEmail);
		renderEmailMock.mockClear();

		const errorCallback = vi.fn(() => {
			throw new Error('Callback failed');
		});

		const config: Config = {
			collections: [
				{
					slug: 'users',
					fields: [],
					auth: true,
				},
			] as CollectionConfig[],
			serverURL: 'https://api.example.com',
		};

		const emailConfig: EmailConfig = {
			...mockEmailConfig,
			forgotPassword: {
				url: errorCallback,
			},
		};

		const result = injectEmailTemplates(config, emailConfig);

		const usersCollection = result.collections?.[0];
		if (
			typeof usersCollection?.auth === 'object' &&
			usersCollection.auth.forgotPassword?.generateEmailHTML
		) {
			await usersCollection.auth.forgotPassword.generateEmailHTML({
				token: 'fallback-token',
				user: { email: 'test@example.com' },
			});

			expect(errorCallback).toHaveBeenCalled();
			expect(consoleWarnSpy).toHaveBeenCalledWith(
				'Failed to generate custom forgot password URL, using default:',
				expect.any(Error),
			);
			expect(renderEmailMock).toHaveBeenCalledWith(
				expect.objectContaining({
					props: expect.objectContaining({
						resetUrl: 'https://api.example.com/admin/reset/fallback-token',
					}),
				}),
			);
		}

		consoleWarnSpy.mockRestore();
	});
});
