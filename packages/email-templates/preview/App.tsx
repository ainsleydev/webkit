// biome-ignore lint/style/useImportType: React is needed for JSX
import * as React from 'react';
import { render } from '@react-email/render';
import { defaultTheme } from '../src/theme/default.js';
import { WelcomeEmail } from './examples/Welcome.js';
import { PasswordResetEmail } from './examples/PasswordReset.js';

type TemplateType = 'welcome' | 'password-reset';

export const App = () => {
	const [selectedTemplate, setSelectedTemplate] = React.useState<TemplateType>('welcome');
	const [html, setHtml] = React.useState<string>('');

	React.useEffect(() => {
		const renderTemplate = () => {
			let component: React.ReactElement;

			switch (selectedTemplate) {
				case 'welcome':
					component = <WelcomeEmail theme={defaultTheme} userName='John Smith' />;
					break;
				case 'password-reset':
					component = <PasswordResetEmail theme={defaultTheme} userName='John Smith' />;
					break;
				default:
					component = <WelcomeEmail theme={defaultTheme} />;
			}

			const rendered = render(component);
			setHtml(rendered);
		};

		renderTemplate();
	}, [selectedTemplate]);

	const copyHtml = () => {
		navigator.clipboard.writeText(html);
		alert('HTML copied to clipboard!');
	};

	return (
		<div
			style={{
				fontFamily: 'system-ui, -apple-system, sans-serif',
				height: '100vh',
				display: 'flex',
				flexDirection: 'column',
			}}
		>
			{/* Header */}
			<div
				style={{
					backgroundColor: '#1a1a1a',
					color: '#fff',
					padding: '20px',
					borderBottom: '1px solid #333',
				}}
			>
				<h1 style={{ margin: '0 0 10px 0', fontSize: '24px' }}>Email Templates Preview</h1>
				<p style={{ margin: '0', color: '#999', fontSize: '14px' }}>
					Preview and test email templates
				</p>
			</div>

			{/* Controls */}
			<div
				style={{
					backgroundColor: '#252525',
					padding: '15px 20px',
					borderBottom: '1px solid #333',
					display: 'flex',
					gap: '10px',
					alignItems: 'centre',
				}}
			>
				<label style={{ color: '#fff', fontWeight: 'bold' }}>Template:</label>
				<select
					value={selectedTemplate}
					onChange={(e) => setSelectedTemplate(e.target.value as TemplateType)}
					style={{
						padding: '8px 12px',
						borderRadius: '4px',
						border: '1px solid #444',
						backgroundColor: '#333',
						color: '#fff',
						cursor: 'pointer',
					}}
				>
					<option value='welcome'>Welcome Email</option>
					<option value='password-reset'>Password Reset</option>
				</select>

				<button
					type='button'
					onClick={copyHtml}
					style={{
						marginLeft: 'auto',
						padding: '8px 16px',
						borderRadius: '4px',
						border: 'none',
						backgroundColor: '#ff5043',
						color: '#fff',
						cursor: 'pointer',
						fontWeight: 'bold',
					}}
				>
					Copy HTML
				</button>
			</div>

			{/* Preview */}
			<div
				style={{
					flex: 1,
					overflow: 'auto',
					backgroundColor: '#1a1a1a',
					padding: '20px',
				}}
			>
				<div
					style={{
						maxWidth: '600px',
						margin: '0 auto',
						backgroundColor: '#fff',
						boxShadow: '0 4px 6px rgba(0, 0, 0, 0.3)',
					}}
				>
					<iframe
						title='Email Preview'
						srcDoc={html}
						style={{
							width: '100%',
							height: '600px',
							border: 'none',
							display: 'block',
						}}
					/>
				</div>
			</div>
		</div>
	);
};
