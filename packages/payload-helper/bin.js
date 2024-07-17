#!/usr/bin/env node
const start = async () => {
	const { bin } = await import('./dist/cli/bin.js');
	await bin();
};

void start();
