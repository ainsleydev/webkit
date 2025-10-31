#!/usr/bin/env node

import { previewCommand } from './preview.js';

/**
 * Main CLI entry point for @ainsleydev/email-templates.
 */
async function main() {
	const args = process.argv.slice(2);
	const command = args[0];

	if (!command || command === 'preview') {
		const directory = args[1] || '.';
		const portArg = args.find((arg) => arg.startsWith('--port='));
		const port = portArg ? Number.parseInt(portArg.split('=')[1], 10) : 3000;

		await previewCommand({ directory, port });
	} else {
		console.error(`Unknown command: ${command}`);
		console.log('\nUsage:');
		console.log('  email-templates preview [directory] [--port=3000]');
		process.exit(1);
	}
}

main().catch((error) => {
	console.error('Error:', error.message);
	process.exit(1);
});
