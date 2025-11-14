#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import { Command } from 'commander';
import { previewEmails } from './commands/preview-emails.js';
import { fileURLToPath } from 'node:url';
import { getPayload, type Payload } from 'payload';

const program = new Command();
const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

program.command('generate-types').description('Generate JSON schema types for Payload CMS');
program
	.command('preview-emails')
	.description('Preview email templates with your Payload configuration')
	.option('-p, --port <number>', 'Port to run preview server on', '3000')
	.action(async (options) => {
		const payload = await getPayloadInstance();

		await previewEmails({
			payload,
			port: Number.parseInt(options.port, 10),
		});
	});

export const bin = async () => {
	await program.parseAsync(process.argv);
};

const getPayloadInstance = async (): Promise<Payload> => {
	const configPath = path.join(process.cwd(), 'src/payload.config.ts');
	const config = (await import(configPath)).default;

	return await getPayload({ config });
};
