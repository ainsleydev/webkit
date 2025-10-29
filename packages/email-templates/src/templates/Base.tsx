import {
	Body,
	Column,
	Container,
	Head,
	Hr,
	Html,
	Img,
	Link,
	Preview,
	Row,
	Section,
	Text,
} from '@react-email/components';
// biome-ignore lint/style/useImportType: React is needed for JSX
import * as React from 'react';
import { generateStyles } from '../theme/styles.js';
import type { EmailTheme } from '../theme/types.js';

interface BaseEmailProps {
	theme: EmailTheme;
	previewText?: string;
	children: React.ReactNode;
}

/**
 * Base email layout component.
 * Provides consistent structure with logo, content area, and footer.
 */
export const BaseEmail = ({ theme, previewText, children }: BaseEmailProps) => {
	const styles = generateStyles(theme);
	const { branding } = theme;

	return (
		<Html>
			<Head />
			{previewText && <Preview>{previewText}</Preview>}
			<Body style={styles.main}>
				<Section style={{ padding: '0 15px' }}>
					<Container style={styles.container}>
						<Section style={styles.logoSection}>
							<Img
								src={branding.logoUrl}
								width={branding.logoWidth.toString()}
								alt={branding.companyName}
							/>
						</Section>
						<Hr style={styles.hr} />
						<Section>{children}</Section>
						<Hr style={styles.hr} />
						<Section style={{ marginTop: '20px' }}>
							<Row>
								<Column align='left'>
									<Text style={styles.footerText}>
										© {new Date().getFullYear()} {branding.companyName}.{' '}
										{branding.footerText}
									</Text>
								</Column>
								{branding.websiteUrl && (
									<Column align='right'>
										<Text style={styles.footerText}>
											<Link
												href={branding.websiteUrl}
												style={{ color: theme.colours.text.action }}
											>
												{new URL(branding.websiteUrl).hostname}
											</Link>
										</Text>
									</Column>
								)}
							</Row>
						</Section>
					</Container>
				</Section>
			</Body>
		</Html>
	);
};
