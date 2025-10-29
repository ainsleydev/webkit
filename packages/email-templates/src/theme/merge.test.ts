import { describe, expect, test } from 'vitest';
import { defaultTheme } from './default.js';
import { mergeTheme } from './merge.js';

describe('mergeTheme', () => {
	test('returns default theme when no partial provided', () => {
		const result = mergeTheme();
		expect(result).toEqual(defaultTheme);
	});

	test('returns default theme when undefined partial provided', () => {
		const result = mergeTheme(undefined);
		expect(result).toEqual(defaultTheme);
	});

	test('merges partial branding with defaults', () => {
		const result = mergeTheme({
			branding: {
				companyName: 'Test Company',
				logoUrl: 'https://test.com/logo.png',
			},
		});

		expect(result.branding.companyName).toBe('Test Company');
		expect(result.branding.logoUrl).toBe('https://test.com/logo.png');
		expect(result.branding.logoWidth).toBe(defaultTheme.branding.logoWidth);
	});

	test('merges partial colours with defaults', () => {
		const result = mergeTheme({
			colours: {
				text: {
					heading: '#000000',
				},
			},
		});

		expect(result.colours.text.heading).toBe('#000000');
		expect(result.colours.text.body).toBe(defaultTheme.colours.text.body);
		expect(result.colours.background).toEqual(defaultTheme.colours.background);
	});

	test('merges multiple partial overrides', () => {
		const result = mergeTheme({
			branding: {
				companyName: 'Multi Test',
			},
			colours: {
				text: {
					action: '#00ff00',
				},
				background: {
					red: '#ff00ff',
				},
			},
		});

		expect(result.branding.companyName).toBe('Multi Test');
		expect(result.branding.logoUrl).toBe(defaultTheme.branding.logoUrl);
		expect(result.colours.text.action).toBe('#00ff00');
		expect(result.colours.text.heading).toBe(defaultTheme.colours.text.heading);
		expect(result.colours.background.red).toBe('#ff00ff');
		expect(result.colours.background.white).toBe(defaultTheme.colours.background.white);
	});
});
