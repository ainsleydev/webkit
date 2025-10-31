import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { previewCommand } from '../src/cli/preview.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Use the CLI preview command with the examples directory.
const examplesDir = join(__dirname, 'examples');

previewCommand({
	directory: examplesDir,
	port: 3000,
}).catch((error) => {
	console.error('Failed to start preview server:', error);
	process.exit(1);
});
