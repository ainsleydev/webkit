#!/usr/bin/env node

import { Command } from 'commander';
import { renderEmail } from './renderer.js';
import type { TemplateName } from './templates/index.js';
import type { PartialEmailTheme } from './theme/types.js';

/**
 * CLI for rendering email templates.
 *
 * Go Integration Pattern:
 * -----------------------
 * This CLI can be called from Go using exec.Command for cross-language email rendering.
 * The CLI accepts JSON input and outputs HTML to stdout, making it easy to integrate.
 *
 * Example Go usage:
 * ```go
 * cmd := exec.Command("npx", "@ainsleydev/email-templates", "render",
 *     "--template", "forgot-password",
 *     "--props", `{"user":{"firstName":"John"},"resetUrl":"https://example.com/reset"}`,
 *     "--theme", `{"branding":{"companyName":"My Company"}}`)
 *
 * output, err := cmd.Output()
 * if err != nil {
 *     return err
 * }
 * html := string(output)
 * ```
 *
 * Benefits:
 * - No runtime Node.js dependency in Go (npx handles execution)
 * - Type-safe on JavaScript side, flexible on Go side
 * - Easy to test independently
 * - Shared templates across JavaScript and Go projects
 */

const program = new Command();

program.name('email-templates').description('Render email templates to HTML').version('0.0.1');

program
	.command('render')
	.description('Render an email template to HTML')
	.requiredOption('-t, --template <name>', 'Template name (e.g., forgot-password)')
	.requiredOption('-p, --props <json>', 'Template props as JSON string')
	.option('-T, --theme <json>', 'Theme overrides as JSON string')
	.option('--plain-text', 'Render as plain text instead of HTML')
	.action(async (options) => {
		try {
			// Parse JSON inputs.
			const props = JSON.parse(options.props);
			const theme: PartialEmailTheme | undefined = options.theme
				? JSON.parse(options.theme)
				: undefined;

			// Render the template.
			const html = await renderEmail({
				template: options.template as TemplateName,
				props,
				theme,
				plainText: options.plainText || false,
			});

			// Output to stdout for Go to capture.
			process.stdout.write(html);
		} catch (error) {
			if (error instanceof Error) {
				process.stderr.write(`Error: ${error.message}\n`);
			}
			process.exit(1);
		}
	});

program.parse();
