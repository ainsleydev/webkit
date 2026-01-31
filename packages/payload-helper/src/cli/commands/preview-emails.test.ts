import { existsSync } from 'node:fs';
import http from 'node:http';
import { resolve } from 'node:path';
import { afterEach, beforeEach, describe, expect, test, vi } from 'vitest';

// Mock node:fs
vi.mock('node:fs', () => ({
	existsSync: vi.fn(),
}));

// Mock @ainsleydev/email-templates
vi.mock('@ainsleydev/email-templates', () => ({
	renderEmail: vi.fn().mockResolvedValue('<html><body>Rendered Email</body></html>'),
}));

// Mock chalk to passthrough
vi.mock('chalk', () => ({
	default: {
		blue: (s: string) => s,
		green: (s: string) => s,
		yellow: (s: string) => s,
		red: (s: string) => s,
		cyan: (s: string) => s,
	},
}));

// Mock email components
vi.mock('../../email/ForgotPasswordEmail.js', () => ({
	ForgotPasswordEmail: () => null,
}));
vi.mock('../../email/VerifyAccountEmail.js', () => ({
	VerifyAccountEmail: () => null,
}));

import { renderEmail } from '@ainsleydev/email-templates';
import { previewEmails } from './preview-emails.js';

const mockedExistsSync = vi.mocked(existsSync);
const mockedRenderEmail = vi.mocked(renderEmail);

describe('previewEmails', () => {
	let server: http.Server | undefined;

	afterEach(() => {
		vi.restoreAllMocks();
		if (server) {
			server.close();
			server = undefined;
		}
	});

	const fetch = async (port: number, path: string): Promise<{ status: number; body: string }> => {
		return new Promise((resolve, reject) => {
			http.get(`http://localhost:${port}${path}`, (res) => {
				let data = '';
				res.on('data', (chunk) => {
					data += chunk;
				});
				res.on('end', () => resolve({ status: res.statusCode || 0, body: data }));
			}).on('error', reject);
		});
	};

	test('should start server and serve index page', async () => {
		const port = 9100;
		mockedExistsSync.mockReturnValue(false);

		// Start server (runs in background)
		const promise = previewEmails({ port });

		// Wait for server to start
		await new Promise((r) => setTimeout(r, 200));

		const res = await fetch(port, '/');
		expect(res.status).toBe(200);
		expect(res.body).toContain('Email Previews');
		expect(res.body).toContain('forgot-password');
		expect(res.body).toContain('verify-account');
	});

	test('should render forgot-password template', async () => {
		const port = 9101;
		mockedExistsSync.mockReturnValue(false);
		mockedRenderEmail.mockResolvedValue('<html>Forgot Password</html>');

		const promise = previewEmails({ port });
		await new Promise((r) => setTimeout(r, 200));

		const res = await fetch(port, '/forgot-password');
		expect(res.status).toBe(200);
		expect(res.body).toContain('Forgot Password');
		expect(mockedRenderEmail).toHaveBeenCalled();
	});

	test('should render verify-account template', async () => {
		const port = 9102;
		mockedExistsSync.mockReturnValue(false);
		mockedRenderEmail.mockResolvedValue('<html>Verify Account</html>');

		const promise = previewEmails({ port });
		await new Promise((r) => setTimeout(r, 200));

		const res = await fetch(port, '/verify-account');
		expect(res.status).toBe(200);
		expect(res.body).toContain('Verify Account');
		expect(mockedRenderEmail).toHaveBeenCalled();
	});

	test('should return 404 for unknown routes', async () => {
		const port = 9103;
		mockedExistsSync.mockReturnValue(false);

		const promise = previewEmails({ port });
		await new Promise((r) => setTimeout(r, 200));

		const res = await fetch(port, '/unknown');
		expect(res.status).toBe(404);
	});

	test('should load config from explicit path', async () => {
		const port = 9104;
		const configPath = 'custom/email.config.ts';
		const resolvedPath = resolve(process.cwd(), configPath);

		mockedExistsSync.mockImplementation((p) => p === resolvedPath);

		// Mock dynamic import
		vi.spyOn(global, 'Function').mockImplementation(() => () => ({}));

		// Since we can't easily mock dynamic import(), just verify it resolves the path
		// The actual import will fail in test, but we verify the flow
		try {
			await previewEmails({ port, configPath });
		} catch {
			// Expected - dynamic import won't work in test
		}

		expect(mockedExistsSync).toHaveBeenCalledWith(resolvedPath);
	});

	test('should auto-detect config files', async () => {
		const port = 9105;

		// Simulate src/email.config.ts existing
		mockedExistsSync.mockImplementation((p) => {
			return p === resolve(process.cwd(), 'src/email.config.ts');
		});

		// Will fail on import but we can verify the detection
		try {
			await previewEmails({ port });
		} catch {
			// Expected - dynamic import won't work in test
		}

		// Should have checked for config files
		expect(mockedExistsSync).toHaveBeenCalledWith(
			resolve(process.cwd(), 'src/email.config.ts'),
		);
	});

	test('should handle render errors gracefully', async () => {
		const port = 9106;
		mockedExistsSync.mockReturnValue(false);
		mockedRenderEmail.mockRejectedValue(new Error('Render failed'));

		const promise = previewEmails({ port });
		await new Promise((r) => setTimeout(r, 200));

		const res = await fetch(port, '/forgot-password');
		expect(res.status).toBe(500);
		expect(res.body).toContain('Render failed');
	});
});
