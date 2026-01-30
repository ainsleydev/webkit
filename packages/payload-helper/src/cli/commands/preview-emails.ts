import { createServer } from 'node:http';
import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import { renderEmail } from '@ainsleydev/email-templates';
import chalk from 'chalk';
import { ForgotPasswordEmail } from '../../email/ForgotPasswordEmail.js';
import { VerifyAccountEmail } from '../../email/VerifyAccountEmail.js';
import type { EmailConfig } from '../../types.js';

/**
 * Default config file names to search for (in order of priority).
 */
const CONFIG_FILE_NAMES = [
	'src/email.config.ts',
	'email.config.ts',
	'src/email.config.js',
	'email.config.js',
];

/**
 * Attempts to resolve the email config file path.
 * If an explicit path is provided, it is used directly.
 * Otherwise, searches for common config file names in the current working directory.
 */
const resolveConfigPath = (explicitPath?: string): string | undefined => {
	if (explicitPath) {
		const resolved = resolve(process.cwd(), explicitPath);
		if (existsSync(resolved)) {
			return resolved;
		}
		console.error(chalk.red(`Config file not found: ${explicitPath}`));
		return undefined;
	}

	for (const name of CONFIG_FILE_NAMES) {
		const candidate = resolve(process.cwd(), name);
		if (existsSync(candidate)) {
			return candidate;
		}
	}

	return undefined;
};

/**
 * Loads the email config from the resolved file path.
 */
const loadEmailConfig = async (configPath: string): Promise<EmailConfig> => {
	const module = await import(configPath);
	return module.default ?? module;
};

/**
 * Available email template previews.
 */
const templates = [
	{ slug: 'forgot-password', label: 'Forgot Password' },
	{ slug: 'verify-account', label: 'Verify Account' },
];

/**
 * Renders a specific email template to HTML.
 */
const renderTemplate = async (slug: string, emailConfig: EmailConfig): Promise<string | null> => {
	const frontEndUrl = emailConfig.frontEndUrl || 'https://yoursite.com';
	const mockUser = { firstName: 'John', email: 'john@example.com' };

	switch (slug) {
		case 'forgot-password':
			return renderEmail({
				component: ForgotPasswordEmail,
				props: {
					user: mockUser,
					resetUrl: `${frontEndUrl}/admin/reset/token123`,
					content: emailConfig.forgotPassword,
				},
				theme: emailConfig.theme,
			});
		case 'verify-account':
			return renderEmail({
				component: VerifyAccountEmail,
				props: {
					user: mockUser,
					verifyUrl: `${frontEndUrl}/admin/verify/token123`,
					content: emailConfig.verifyAccount,
				},
				theme: emailConfig.theme,
			});
		default:
			return null;
	}
};

/**
 * Generates the index HTML page listing available templates.
 */
const renderIndexPage = (): string => {
	const links = templates
		.map((t) => `<li style="margin-bottom: 8px;"><a href="/${t.slug}">${t.label}</a></li>`)
		.join('\n');

	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Email Previews</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 40px auto; padding: 0 20px; color: #333; }
		h1 { margin-bottom: 24px; }
		ul { list-style: none; padding: 0; }
		a { color: #2563eb; text-decoration: none; font-size: 18px; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<h1>Email Previews</h1>
	<ul>${links}</ul>
</body>
</html>`;
};

export const previewEmails = async (options: { configPath?: string; port?: number }) => {
	const port = options.port || 3000;

	console.log(chalk.blue('Looking for email configuration...'));

	const configPath = resolveConfigPath(options.configPath);

	let emailConfig: EmailConfig = {};

	if (!configPath) {
		console.log(chalk.yellow('No email config file found.'));
		console.log(chalk.yellow(`Searched for: ${CONFIG_FILE_NAMES.join(', ')}`));
		console.log(chalk.yellow('\nCreate an email config file, for example:\n'));
		console.log(
			chalk.cyan(`  // src/email.config.ts
  import { defineEmailConfig } from '@ainsleydev/payload-helper';

  export default defineEmailConfig({
    frontEndUrl: 'https://yoursite.com',
    theme: { /* ... */ },
  });`),
		);
		console.log(chalk.yellow('\nUsing default theme for previews.\n'));
	} else {
		console.log(chalk.green(`Found config: ${configPath}`));
		emailConfig = await loadEmailConfig(configPath);
	}

	const server = createServer(async (req, res) => {
		const url = new URL(req.url || '/', `http://localhost:${port}`);
		const pathname = url.pathname.replace(/^\/+|\/+$/g, '');

		// Index page
		if (!pathname) {
			res.writeHead(200, { 'Content-Type': 'text/html' });
			res.end(renderIndexPage());
			return;
		}

		// Template page
		try {
			const html = await renderTemplate(pathname, emailConfig);
			if (html) {
				res.writeHead(200, { 'Content-Type': 'text/html' });
				res.end(html);
				return;
			}
		} catch (err) {
			console.error(chalk.red(`Error rendering template "${pathname}":`), err);
			res.writeHead(500, { 'Content-Type': 'text/plain' });
			res.end(`Error rendering template: ${err}`);
			return;
		}

		// 404
		res.writeHead(404, { 'Content-Type': 'text/html' });
		res.end(renderIndexPage());
	});

	server.listen(port, () => {
		console.log(chalk.green(`\nEmail preview server running at http://localhost:${port}\n`));
		for (const t of templates) {
			console.log(chalk.cyan(`  ${t.label}: http://localhost:${port}/${t.slug}`));
		}
		console.log('');
	});

	// Graceful shutdown
	const shutdown = () => {
		console.log(chalk.blue('\nShutting down preview server...'));
		server.close(() => {
			process.exit(0);
		});
	};

	process.on('SIGINT', shutdown);
	process.on('SIGTERM', shutdown);
};
