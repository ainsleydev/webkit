import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Preprocess Svelte components with TypeScript support
	preprocess: vitePreprocess(),

	compilerOptions: {
		runes: true,
	},

	kit: {
		// No adapter needed for library packaging
	},
};

export default config;
