#!/usr/bin/env node
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { Command } from 'commander';
import { previewEmails } from './commands/preview-emails.js';

const program = new Command();
const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

program.command('generate-types').description('Generate JSON schema types for Payload CMS');
program
	.command('preview-emails')
	.description('Preview email templates with your email configuration')
	.option('-p, --port <number>', 'Port to run preview server on', '3000')
	.option('-c, --config <path>', 'Path to email config file')
	.action(async (options) => {
		await previewEmails({
			configPath: options.config,
			port: Number.parseInt(options.port, 10),
		});
	});

export const bin = async () => {
	await program.parseAsync(process.argv);
};
