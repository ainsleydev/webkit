#!/usr/bin/env node
import chalk from 'chalk';
import { Command } from 'commander';
import { previewEmails } from './commands/preview-emails.js';

const program = new Command();

program.command('generate-types').description('Generate JSON schema types for Payload CMS');

program
	.command('preview-emails')
	.description('Preview email templates with your Payload configuration')
	.option('-p, --port <number>', 'Port to run preview server on', '3000')
	.action(async (options) => {
		await previewEmails({
			port: Number.parseInt(options.port, 10),
		});
	});

export const bin = async () => {
	await program.parseAsync(process.argv);
};
