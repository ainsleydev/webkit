#!/usr/bin/env node
import { Command } from 'commander';
import { generateTypes } from './types.js';
import chalk from 'chalk';

const program = new Command();

program
	.command('generate-types')
	.description('Generate JSON schema types for Payload CMS')
	.action(async () => {
		try {
			await generateTypes();
		} catch (error) {
			console.log(chalk.red(error));
		}
	});

export const bin = async () => {
	await program.parseAsync(process.argv);
};
