// biome-ignore lint/style/useImportType: React is needed for JSX
import * as React from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App.js';

const rootElement = document.getElementById('root');
if (!rootElement) {
	throw new Error('Root element not found');
}

const root = createRoot(rootElement);
root.render(<App />);
