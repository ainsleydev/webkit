import { createServer } from 'node:http';
import { readdir, stat } from 'node:fs/promises';
import { join, relative, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
import { render } from '@react-email/render';
import * as React from 'react';
import { defaultTheme } from '../theme/default.js';
import type { EmailTheme } from '../theme/types.js';

type EmailComponent = React.ComponentType<{ theme: EmailTheme }>;

interface PreviewOptions {
	directory: string;
	port: number;
}

interface TemplateInfo {
	path: string;
	name: string;
	route: string;
}

/**
 * Discovers email template files in a directory.
 *
 * @param directory - The directory to search for templates
 * @returns Array of discovered template information
 */
async function discoverTemplates(directory: string): Promise<TemplateInfo[]> {
	const templates: TemplateInfo[] = [];
	const resolvedDir = resolve(directory);

	try {
		const entries = await readdir(resolvedDir);

		for (const entry of entries) {
			const fullPath = join(resolvedDir, entry);
			const stats = await stat(fullPath);

			if (stats.isFile() && (entry.endsWith('.tsx') || entry.endsWith('.jsx'))) {
				const name = entry.replace(/\.(tsx|jsx)$/, '');
				const route = `/${name.toLowerCase()}`;

				templates.push({
					path: fullPath,
					name,
					route,
				});
			}
		}
	} catch (error) {
		const message = error instanceof Error ? error.message : String(error);
		throw new Error(`Failed to read directory "${directory}": ${message}`);
	}

	return templates;
}

/**
 * Dynamically imports a template component.
 *
 * @param templatePath - Absolute path to the template file
 * @returns The imported component
 */
async function loadTemplate(templatePath: string): Promise<EmailComponent> {
	try {
		// Convert file path to file URL for dynamic import.
		const fileUrl = pathToFileURL(templatePath).href;
		const module = await import(fileUrl);

		// Try to find the component - check default export first, then named exports.
		const component =
			module.default || Object.values(module).find((exp) => typeof exp === 'function');

		if (!component) {
			throw new Error('No React component found in template file');
		}

		return component as EmailComponent;
	} catch (error) {
		const message = error instanceof Error ? error.message : String(error);
		throw new Error(`Failed to load template from "${templatePath}": ${message}`);
	}
}

/**
 * Creates and starts the preview server.
 *
 * @param options - Preview server options
 */
export async function previewCommand(options: PreviewOptions): Promise<void> {
	const { directory, port } = options;

	console.log('\nDiscovering email templates...');
	const templates = await discoverTemplates(directory);

	if (templates.length === 0) {
		console.error(`\nNo email templates found in "${directory}"`);
		console.log('Expected .tsx or .jsx files with React components.');
		process.exit(1);
	}

	console.log(`Found ${templates.length} template(s):\n`);
	for (const template of templates) {
		console.log(`  - ${template.name}`);
	}

	// Pre-load all templates.
	const templateMap = new Map<string, EmailComponent>();

	for (const template of templates) {
		try {
			const component = await loadTemplate(template.path);
			templateMap.set(template.route, component);
		} catch (error) {
			const message = error instanceof Error ? error.message : String(error);
			console.error(`\nWarning: Failed to load ${template.name}: ${message}`);
		}
	}

	// Create HTTP server.
	const server = createServer(async (req, res) => {
		const url = req.url || '/';

		// Handle root - redirect to first template.
		if (url === '/') {
			const firstTemplate = templates[0];
			if (firstTemplate) {
				res.writeHead(302, { Location: firstTemplate.route });
				res.end();
				return;
			}
		}

		// Try to render the requested template.
		const component = templateMap.get(url);

		if (component) {
			try {
				// Create element with default theme.
				const element = React.createElement(component, { theme: defaultTheme });
				const html = await render(element);

				res.writeHead(200, { 'Content-Type': 'text/html' });
				res.end(html);
				return;
			} catch (error) {
				const message = error instanceof Error ? error.message : String(error);
				res.writeHead(500, { 'Content-Type': 'text/plain' });
				res.end(`Error rendering template: ${message}`);
				return;
			}
		}

		// 404 - template not found.
		res.writeHead(404, { 'Content-Type': 'text/html' });
		res.end(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>404 - Not Found</title>
					<style>
						body { font-family: system-ui; padding: 40px; background: #1a1a1a; color: #fff; }
						h1 { color: #ff5043; }
						a { color: #ff5043; text-decoration: none; }
						a:hover { text-decoration: underline; }
						ul { list-style: none; padding: 0; }
						li { margin: 8px 0; }
					</style>
				</head>
				<body>
					<h1>404 - Template not found</h1>
					<p>Available templates:</p>
					<ul>
						${templates.map((t) => `<li><a href="${t.route}">${t.name}</a></li>`).join('\n')}
					</ul>
				</body>
			</html>
		`);
	});

	// Start server.
	server.listen(port, () => {
		console.log('\n  Email Templates Preview');
		console.log(`  ➜  Local:   http://localhost:${port}/\n`);
		console.log('  Available templates:');
		for (const template of templates) {
			console.log(`  - http://localhost:${port}${template.route}`);
		}
		console.log('\n  Press Ctrl+C to stop\n');
	});
}
