#!/usr/bin/env node

/**
 * CLI entry point that uses tsx to handle TypeScript/JSX files.
 */

import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { dirname, join } from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Run the TypeScript source directly with tsx.
const cliPath = join(__dirname, 'src', 'cli', 'index.ts');

const child = spawn('npx', ['tsx', cliPath, ...process.argv.slice(2)], {
	stdio: 'inherit',
	shell: true,
});

child.on('exit', (code) => {
	process.exit(code || 0);
});
