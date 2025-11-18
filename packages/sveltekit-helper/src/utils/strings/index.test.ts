import { describe, expect, test } from 'vitest';

import { generateRandomString } from './index';

describe('generateRandomString', () => {
	test('generates a string of the specified length', () => {
		const length = 10;
		const result = generateRandomString(length);
		expect(result).toHaveLength(length);
	});

	test('generates alphanumeric characters only', () => {
		const result = generateRandomString(100);
		expect(result).toMatch(/^[a-z0-9]+$/);
	});

	test('generates different strings on multiple calls', () => {
		const first = generateRandomString(20);
		const second = generateRandomString(20);
		expect(first).not.toBe(second);
	});

	test('handles length of 0', () => {
		const result = generateRandomString(0);
		expect(result).toBe('');
	});

	test('handles length of 1', () => {
		const result = generateRandomString(1);
		expect(result).toHaveLength(1);
	});
});
