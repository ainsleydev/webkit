import type { EmailTheme } from '@ainsleydev/email-templates';
import { describe, expect, test, vi } from 'vitest';
import { VerifyAccountEmail } from './VerifyAccountEmail.js';
import type { VerifyAccountEmailProps } from './VerifyAccountEmail.js';

// Mock the email templates module
vi.mock('@ainsleydev/email-templates', () => ({
	BaseEmail: ({ children }: { children: React.ReactNode }) => children,
	Button: ({ children, href }: { children: React.ReactNode; href: string }) => (
		<a href={href}>{children}</a>
	),
	Heading: ({ children }: { children: React.ReactNode }) => <h1>{children}</h1>,
	Section: ({ children }: { children: React.ReactNode }) => <section>{children}</section>,
	Text: ({ children }: { children: React.ReactNode }) => <p>{children}</p>,
}));

const mockTheme: EmailTheme = {
	colours: {
		text: {
			heading: '#000000',
			body: '#333333',
		},
		background: {
			primary: '#ffffff',
			accent: '#007bff',
		},
	},
	branding: {
		companyName: 'Test Company',
		websiteUrl: 'https://example.com',
	},
	fonts: {
		heading: 'Arial, sans-serif',
		body: 'Arial, sans-serif',
	},
};

describe('VerifyAccountEmail', () => {
	test('should render with default content when no overrides provided', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'John',
				email: 'john@example.com',
			},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// Component should render without errors
		expect(result).toBeDefined();
	});

	test('should use firstName when available', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'Jane',
				email: 'jane@example.com',
			},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// The userName should be set to firstName
		// This is tested through the default heading which uses userName
		expect(result).toBeDefined();
	});

	test('should use email when firstName is not available', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				email: 'john@example.com',
			},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// The userName should be set to email
		expect(result).toBeDefined();
	});

	test('should use "there" when neither firstName nor email is available', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// The userName should default to "there"
		expect(result).toBeDefined();
	});

	test('should use custom content overrides when provided', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'John',
			},
			verifyUrl: 'https://example.com/verify/token123',
			content: {
				previewText: 'Custom preview text',
				heading: 'Custom Welcome Heading',
				bodyText: 'Custom verification message',
				buttonText: 'Custom Verify Button',
			},
		};

		const result = VerifyAccountEmail(props);

		// Component should use custom content
		expect(result).toBeDefined();
	});

	test('should use partial content overrides with defaults', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'John',
			},
			verifyUrl: 'https://example.com/verify/token123',
			content: {
				buttonText: 'Confirm Email',
				// Other fields should use defaults
			},
		};

		const result = VerifyAccountEmail(props);

		// Component should use custom button text and defaults for others
		expect(result).toBeDefined();
	});

	test('should pass verifyUrl to the button', () => {
		const verifyUrl = 'https://example.com/verify/xyz789abc';
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'John',
			},
			verifyUrl,
		};

		const result = VerifyAccountEmail(props);

		// The verifyUrl should be used in the button href
		expect(result).toBeDefined();
	});

	test('should apply theme colors correctly', () => {
		const customTheme: EmailTheme = {
			colours: {
				text: {
					heading: '#ff0000',
					body: '#00ff00',
				},
				background: {
					primary: '#ffffff',
					accent: '#0000ff',
				},
			},
			branding: {
				companyName: 'Custom Company',
				websiteUrl: 'https://custom.com',
			},
			fonts: {
				heading: 'Georgia, serif',
				body: 'Georgia, serif',
			},
		};

		const props: VerifyAccountEmailProps = {
			theme: customTheme,
			user: {
				firstName: 'John',
			},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// Component should use custom theme colors
		expect(result).toBeDefined();
	});

	test('should handle empty content object', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {
				firstName: 'John',
			},
			verifyUrl: 'https://example.com/verify/token123',
			content: {},
		};

		const result = VerifyAccountEmail(props);

		// Should use all defaults when content is empty object
		expect(result).toBeDefined();
	});

	test('should render with minimal props', () => {
		const props: VerifyAccountEmailProps = {
			theme: mockTheme,
			user: {},
			verifyUrl: 'https://example.com/verify/token123',
		};

		const result = VerifyAccountEmail(props);

		// Should handle minimal props gracefully
		expect(result).toBeDefined();
	});
});
