import { spawn } from 'node:child_process';
import { mkdirSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import chalk from 'chalk';
import type { Payload } from 'payload';
import type { PayloadHelperPluginConfig } from '../../types.js';

/**
 * Escapes a string for safe use in template literals.
 */
const escapeForTemplate = (str: string): string => {
	return str.replace(/\\/g, '\\\\').replace(/`/g, '\\`').replace(/\$/g, '\\$');
};

/**
 * Retrieves the email configuration from the Payload config.
 */
const getEmailConfig = (payload: Payload): PayloadHelperPluginConfig['email'] | undefined => {
	const custom = payload.config.custom as Record<string, unknown> | undefined;
	const payloadHelperOptions = custom?.payloadHelperOptions as
		| PayloadHelperPluginConfig
		| undefined;
	return payloadHelperOptions?.email;
};

export const previewEmails = async (options: { payload: Payload; port?: number }) => {
	const port = options.port || 3000;
	const payload = options.payload;

	console.log(chalk.blue('🔍 Looking for email configuration...'));

	// Get email config from stored plugin options
	const emailConfig = getEmailConfig(payload);

	if (!emailConfig) {
		console.log(chalk.yellow('⚠️  No email configuration found'));
		console.log(
			chalk.yellow('   Make sure you have configured email in your payloadHelper plugin:\n'),
		);
		console.log(
			chalk.cyan(`   payloadHelper({
     email: {
       theme: { /* ... */ },
       frontEndUrl: 'https://yoursite.com',
     }
   })`),
		);
		console.log(chalk.yellow('\n   Using default theme for email previews'));
	} else {
		console.log(chalk.green('✓ Found email configuration'));
	}

	// Create temp directory for preview files
	const tempDir = join(tmpdir(), `payload-helper-preview-${Date.now()}-${process.pid}`);
	mkdirSync(tempDir, { recursive: true });

	console.log(chalk.blue('📝 Generating preview templates...'));

	// Extract theme and frontEndUrl
	const themeConfig = emailConfig?.theme ? JSON.stringify(emailConfig.theme, null, 2) : '{}';
	const frontEndUrl = emailConfig?.frontEndUrl || 'https://yoursite.com';

	// Safely escape the frontEndUrl for use in template strings
	const escapedFrontEndUrl = escapeForTemplate(frontEndUrl);

	// Generate ForgotPassword preview
	const forgotPasswordContent = emailConfig?.forgotPassword
		? JSON.stringify(emailConfig.forgotPassword, null, 3)
		: 'undefined';

	const forgotPasswordPreview = `import { renderEmail } from '@ainsleydev/email-templates';
import { ForgotPasswordEmail } from '@ainsleydev/payload-helper';

export default async function render() {
	return renderEmail({
		component: ForgotPasswordEmail,
		props: {
			user: { firstName: 'John', email: 'john@example.com' },
			resetUrl: \`${escapedFrontEndUrl}/admin/reset/token123\`,
			content: ${forgotPasswordContent},
		},
		theme: ${themeConfig},
	});
}
`;

	// Generate VerifyAccount preview
	const verifyAccountContent = emailConfig?.verifyAccount
		? JSON.stringify(emailConfig.verifyAccount, null, 3)
		: 'undefined';

	const verifyAccountPreview = `import { renderEmail } from '@ainsleydev/email-templates';
import { VerifyAccountEmail } from '@ainsleydev/payload-helper';

export default async function render() {
	return renderEmail({
		component: VerifyAccountEmail,
		props: {
			user: { firstName: 'John', email: 'john@example.com' },
			verifyUrl: \`${escapedFrontEndUrl}/admin/verify/token123\`,
			content: ${verifyAccountContent},
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
