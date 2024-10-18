#!/usr/bin/env node
import chalk from 'chalk';
import { Command } from 'commander';

const program = new Command();

program
	.command('generate-types')
	.description('Generate JSON schema types for Payload CMS')


export const bin = async () => {
	await program.parseAsync(process.argv);
};
