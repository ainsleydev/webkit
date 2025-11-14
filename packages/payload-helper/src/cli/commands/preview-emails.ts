import chalk from 'chalk';
import { spawn } from 'node:child_process';
import { existsSync, mkdirSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';
import type { PayloadHelperPluginConfig } from '../../types.js';

export const previewEmails = async (options: { port?: number }) => {
	const port = options.port || 3000;

	console.log(chalk.blue('🔍 Looking for payload.config.ts...'));

	// Find payload.config.ts in current directory
	const configPath = join(process.cwd(), 'payload.config.ts');
	if (!existsSync(configPath)) {
		console.error(chalk.red('❌ Could not find payload.config.ts in current directory'));
		process.exit(1);
	}

	console.log(chalk.green('✓ Found payload.config.ts'));

	// Load the config
	let emailConfig: PayloadHelperPluginConfig['email'];
	try {
		const configUrl = pathToFileURL(configPath).href;
		const configModule = await import(configUrl);
		const config = configModule.default || configModule;

		// Try to find payloadHelper plugin config
		const plugins = config.plugins || [];
		const helperPlugin = plugins.find((p: unknown) => {
			if (typeof p === 'object' && p !== null) {
				const plugin = p as Record<string, unknown>;
				const pluginOptions = plugin.pluginOptions as Record<string, unknown> | undefined;
				const pluginConfig = plugin.config as Record<string, unknown> | undefined;
				return pluginOptions?.email !== undefined || pluginConfig?.email !== undefined;
			}
			return false;
		});

		if (helperPlugin && typeof helperPlugin === 'object') {
			const plugin = helperPlugin as Record<string, unknown>;
			const pluginOptions = plugin.pluginOptions as Record<string, unknown> | undefined;
			const pluginConfig = plugin.config as Record<string, unknown> | undefined;
			emailConfig = (pluginOptions?.email ||
				pluginConfig?.email) as PayloadHelperPluginConfig['email'];
		}

		if (!emailConfig) {
			console.log(chalk.yellow('⚠️  No email configuration found in payload.config.ts'));
			console.log(chalk.yellow('   Using default theme for email previews'));
		} else {
			console.log(chalk.green('✓ Found email configuration'));
		}
	} catch (error) {
		console.error(chalk.red('❌ Error loading payload.config.ts:'), error);
		process.exit(1);
	}

	// Create temp directory for preview files
	const tempDir = join(tmpdir(), `payload-helper-preview-${Date.now()}`);
	mkdirSync(tempDir, { recursive: true });

	console.log(chalk.blue('📝 Generating preview templates...'));

	// Extract theme configuration
	const themeConfig = emailConfig?.theme ? JSON.stringify(emailConfig.theme, null, 2) : '{}';
	const frontEndUrl = emailConfig?.frontEndUrl || 'https://yoursite.com';

	// Generate ForgotPassword preview
	const forgotPasswordPreview = `import { renderEmail } from '@ainsleydev/email-templates';
import { ForgotPasswordEmail } from '@ainsleydev/payload-helper';

export default async function render() {
	return renderEmail({
		component: ForgotPasswordEmail,
		props: {
			user: { firstName: 'John', email: 'john@example.com' },
			resetUrl: '${frontEndUrl}/admin/reset/token123',
			content: ${emailConfig?.forgotPassword ? JSON.stringify(emailConfig.forgotPassword, null, 3) : 'undefined'},
		},
		theme: ${themeConfig},
	});
}
`;

	// Generate VerifyAccount preview
	const verifyAccountPreview = `import { renderEmail } from '@ainsleydev/email-templates';
import { VerifyAccountEmail } from '@ainsleydev/payload-helper';

export default async function render() {
	return renderEmail({
		component: VerifyAccountEmail,
		props: {
			user: { firstName: 'John', email: 'john@example.com' },
			verifyUrl: '${frontEndUrl}/admin/verify/token123',
			content: ${emailConfig?.verifyAccount ? JSON.stringify(emailConfig.verifyAccount, null, 3) : 'undefined'},
		},
		theme: ${themeConfig},
	});
}
`;

	// Write preview files
	writeFileSync(join(tempDir, 'forgot-password-email-preview.tsx'), forgotPasswordPreview);
	writeFileSync(join(tempDir, 'verify-account-email-preview.tsx'), verifyAccountPreview);

	console.log(chalk.green('✓ Preview templates generated'));
	console.log(chalk.blue(`🚀 Starting preview server on http://localhost:${port}...`));

	// Launch email-templates preview
	const previewProcess = spawn('npx', ['email-templates', 'preview', tempDir, `--port=${port}`], {
		stdio: 'inherit',
		shell: true,
	});

	// Cleanup on exit
	const cleanup = () => {
		console.log(chalk.blue('\n🧹 Cleaning up...'));
		try {
			rmSync(tempDir, { recursive: true, force: true });
			console.log(chalk.green('✓ Cleanup complete'));
		} catch (error) {
			console.error(chalk.red('❌ Error during cleanup:'), error);
		}
	};

	previewProcess.on('exit', (code) => {
		cleanup();
		process.exit(code || 0);
	});

	// Handle Ctrl+C
	process.on('SIGINT', () => {
		previewProcess.kill('SIGINT');
	});

	process.on('SIGTERM', () => {
		previewProcess.kill('SIGTERM');
	});
};
