import path from 'node:path';
import process from 'node:process';
import { previewCommand } from '../src/cli/preview.js';

const __dirname = path.resolve(path.dirname(process.argv[1]));
const examplesDir = path.join(__dirname, 'examples');

previewCommand({
	directory: examplesDir,
	port: 3000,
}).catch((error) => {
	console.error('Failed to start preview server:', error);
	process.exit(1);
});
