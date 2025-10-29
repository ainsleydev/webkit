import { createServer } from 'node:http';
import { render } from '@react-email/render';
import React from 'react';
import { defaultTheme } from '../src/theme/default.js';
import { WelcomeEmail } from './examples/Welcome.js';
import { PasswordResetEmail } from './examples/PasswordReset.js';

const PORT = 3000;

const server = createServer(async (req, res) => {
	const url = req.url;

	try {
		let html: string;

		if (url === '/' || url === '/welcome') {
			const template = React.createElement(WelcomeEmail, {
				theme: defaultTheme,
				userName: 'John Smith',
			});
			html = await render(template);
		} else if (url === '/password-reset') {
			const template = React.createElement(PasswordResetEmail, {
				theme: defaultTheme,
				userName: 'John Smith',
			});
			html = await render(template);
		} else {
			res.writeHead(404, { 'Content-Type': 'text/html' });
			res.end(`
				<!DOCTYPE html>
				<html>
					<head>
						<title>404 - Not Found</title>
						<style>
							body { font-family: system-ui; padding: 40px; background: #1a1a1a; color: #fff; }
							h1 { color: #ff5043; }
							a { color: #ff5043; }
						</style>
					</head>
					<body>
						<h1>404 - Template not found</h1>
						<p>Available templates:</p>
						<ul>
							<li><a href="/welcome">Welcome Email</a></li>
							<li><a href="/password-reset">Password Reset Email</a></li>
						</ul>
					</body>
				</html>
			`);
			return;
		}

		res.writeHead(200, { 'Content-Type': 'text/html' });
		res.end(html);
	} catch (error) {
		res.writeHead(500, { 'Content-Type': 'text/plain' });
		res.end(`Error rendering template: ${error.message}`);
	}
});

server.listen(PORT, () => {
	console.log('\n  Email Templates Preview');
	console.log(`  ➜  Local:   http://localhost:${PORT}/`);
	console.log('\n  Available templates:');
	console.log(`  - http://localhost:${PORT}/welcome`);
	console.log(`  - http://localhost:${PORT}/password-reset\n`);
});
