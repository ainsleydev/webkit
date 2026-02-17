import type { Field, FieldHookArgs, TextField, TypeWithID } from 'payload';
import { describe, expect, test, vi } from 'vitest';
import { URLField } from './URLField.js';

describe('URLField', () => {
	test('returns a text field with correct defaults', () => {
		const field = URLField({
			generate: () => 'https://example.com',
		});

		expect(field).toMatchObject({
			name: 'url',
			label: 'URL',
			type: 'text',
			virtual: true,
			admin: {
				readOnly: true,
				position: 'sidebar',
			},
		});
	});

	test('calls generate function in afterRead hook', async () => {
		const generate = vi.fn(() => 'https://example.com/page');

		const field = URLField({ generate });
		const hooks = (field as TextField).hooks?.afterRead;
		expect(hooks).toHaveLength(1);

		const result = await hooks?.[0]({ draft: false } as unknown as FieldHookArgs<TypeWithID>);
		expect(generate).toHaveBeenCalled();
		expect(result).toBe('https://example.com/page');
	});

	test('appends draft query parameter when in draft mode', async () => {
		const field = URLField({
			generate: () => 'https://example.com/page',
		});

		const hooks = (field as TextField).hooks?.afterRead;
		const result = await hooks?.[0]({ draft: true } as unknown as FieldHookArgs<TypeWithID>);
		expect(result).toBe('https://example.com/page?draft=true');
	});

	test('appends draft parameter correctly when URL already has query params', async () => {
		const field = URLField({
			generate: () => 'https://example.com/page?foo=bar',
		});

		const hooks = (field as TextField).hooks?.afterRead;
		const result = await hooks?.[0]({ draft: true } as unknown as FieldHookArgs<TypeWithID>);
		expect(result).toBe('https://example.com/page?foo=bar&draft=true');
	});

	test('handles async generate function', async () => {
		const field = URLField({
			generate: async () => 'https://example.com/async',
		});

		const hooks = (field as TextField).hooks?.afterRead;
		const result = await hooks?.[0]({ draft: false } as unknown as FieldHookArgs<TypeWithID>);
		expect(result).toBe('https://example.com/async');
	});

	test('applies overrides to the base field', () => {
		const field = URLField({
			generate: () => 'https://example.com',
			overrides: {
				name: 'customUrl',
				label: 'Custom URL',
			},
		});

		expect(field).toMatchObject({
			name: 'customUrl',
			label: 'Custom URL',
			type: 'text',
		});
	});

	test('returns undefined when generate returns undefined', async () => {
		const field = URLField({ generate: async () => undefined });
		const hooks = (field as TextField).hooks?.afterRead;
		const result = await hooks?.[0]({ draft: true } as unknown as FieldHookArgs<TypeWithID>);
		expect(result).toBeUndefined();
	});

	test('uses empty object when overrides not provided', () => {
		const field = URLField({
			generate: () => 'https://example.com',
		});

		expect(field).toMatchObject({
			name: 'url',
			type: 'text',
		});
	});
});
